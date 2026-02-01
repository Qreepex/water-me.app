package middlewares

import (
	"log"
	"net/http"
	"sync"
	"time"
)

// RateLimiter tracks requests per identifier
type RateLimiter struct {
	limiter map[string]*userLimiter
	mu      sync.RWMutex
	cleanup *time.Ticker
}

type userLimiter struct {
	count       int
	lastReset   time.Time
	requestsIP  int
	lastResetIP time.Time
}

var (
	// Per-user limit: 100 requests per minute
	userRequestsPerMinute = 100
	// Per-IP limit: 1000 requests per minute
	ipRequestsPerMinute = 1000
)

// NewRateLimiter creates a new rate limiter with automatic cleanup
func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		limiter: make(map[string]*userLimiter),
		cleanup: time.NewTicker(5 * time.Minute),
	}

	// Cleanup old entries periodically
	go func() {
		for range rl.cleanup.C {
			rl.mu.Lock()
			now := time.Now()
			for key, limiter := range rl.limiter {
				// Remove entries that haven't been used in 10 minutes
				if now.Sub(limiter.lastReset) > 10*time.Minute &&
					now.Sub(limiter.lastResetIP) > 10*time.Minute {
					delete(rl.limiter, key)
				}
			}
			rl.mu.Unlock()
		}
	}()

	return rl
}

// getRealIP gets the real client IP, considering Cloudflare headers
func getRealIP(r *http.Request) string {
	// Cloudflare header
	if ip := r.Header.Get("CF-Connecting-IP"); ip != "" {
		return ip
	}
	// Fallback to X-Forwarded-For
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	// Fallback to direct connection
	return r.RemoteAddr
}

// IsRateLimited checks if a user or IP has exceeded the rate limit
func (rl *RateLimiter) IsRateLimited(userID string, ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Get or create limiter for this user-IP combination
	key := userID + "|" + ip
	limiter, exists := rl.limiter[key]
	if !exists {
		limiter = &userLimiter{
			count:       0,
			lastReset:   now,
			requestsIP:  0,
			lastResetIP: now,
		}
		rl.limiter[key] = limiter
	}

	// Check user rate limit (per-user)
	if now.Sub(limiter.lastReset) > time.Minute {
		limiter.count = 0
		limiter.lastReset = now
	}
	limiter.count++
	if limiter.count > userRequestsPerMinute {
		log.Printf("Rate limit exceeded for user %s: %d requests", userID, limiter.count)
		return true
	}

	// Check IP rate limit (per-IP, regardless of user)
	if now.Sub(limiter.lastResetIP) > time.Minute {
		limiter.requestsIP = 0
		limiter.lastResetIP = now
	}
	limiter.requestsIP++
	if limiter.requestsIP > ipRequestsPerMinute {
		log.Printf("Rate limit exceeded for IP %s: %d requests", ip, limiter.requestsIP)
		return true
	}

	return false
}

// RateLimitMiddleware returns a middleware that enforces rate limiting
func RateLimitMiddleware(rl *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getRealIP(r)
			userID := r.Header.Get("X-User-ID") // Set by auth middleware

			if rl.IsRateLimited(userID, ip) {
				w.Header().Set("Retry-After", "60")
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
