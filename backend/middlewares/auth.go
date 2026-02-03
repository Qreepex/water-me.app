package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/qreepex/water-me-app/backend/constants"
	"github.com/qreepex/water-me-app/backend/services"
)

type ctxKey string

func WithUserID(r *http.Request, userID string) *http.Request {
	ctx := context.WithValue(r.Context(), constants.UserIdKey, userID)
	return r.WithContext(ctx)
}

// GetUserID extracts userID from request context.
func GetUserID(r *http.Request) (string, bool) {
	id, ok := r.Context().Value(constants.UserIdKey).(string)
	return id, ok
}

// AuthMiddleware validates Bearer JWT and injects userID into context.
func AuthMiddleware(firebase *services.FirebaseService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			auth(next.ServeHTTP, firebase)(w, r)
		})
	}
}

func auth(next http.HandlerFunc, firebase *services.FirebaseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Allow unauthenticated access to /api/stats
		if r.URL.Path == "/api/stats" {
			next(w, r)
			return
		}

		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		token, err := firebase.VerifyIDToken(r.Context(), tokenStr)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		uid := token.UID

		next(w, WithUserID(r, uid))
	}
}
