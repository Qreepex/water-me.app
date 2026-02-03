package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/qreepex/water-me-app/backend/util"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	jwtSecretEnvKey = "JWT_SECRET"
	jwtExpiration   = 7 * 24 * time.Hour // 7 days
)

type contextKey string

const userIDKey contextKey = "userID"

// HashPassword hashes a plain-text password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(bytes), nil
}

// VerifyPassword checks if a plain-text password matches the hashed password.
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// GenerateJWT creates a JWT token for the given user ID.
func GenerateJWT(userID string, secret string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExpiration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("sign JWT: %w", err)
	}
	return signedToken, nil
}

// VerifyJWT validates a JWT token and returns the user ID.
func VerifyJWT(tokenString string, secret string) (string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		},
	)
	if err != nil {
		return "", fmt.Errorf("parse JWT: %w", err)
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims.Subject, nil
	}

	return "", errors.New("invalid token")
}

// authMiddleware extracts and verifies JWT from Authorization header.
func authMiddleware(secret string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Allow OPTIONS requests through without authentication
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			util.RespondJSON(
				w,
				http.StatusUnauthorized,
				map[string]string{"error": "Missing authorization header"},
			)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			util.RespondJSON(
				w,
				http.StatusUnauthorized,
				map[string]string{"error": "Invalid authorization header format"},
			)
			return
		}

		tokenString := parts[1]
		userID, err := VerifyJWT(tokenString, secret)
		if err != nil {
			util.RespondJSON(
				w,
				http.StatusUnauthorized,
				map[string]string{"error": "Invalid or expired token"},
			)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// getUserID extracts the user ID from the request context.
func getUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(userIDKey).(string)
	return userID, ok
}
