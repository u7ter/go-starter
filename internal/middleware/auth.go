package middleware

import (
	"context"
	"net/http"
	"strings"

	"go-starter/internal/models"
	"go-starter/internal/services"
)

type contextKey string

const userIDKey contextKey = "user_id"

// AuthMiddleware creates a middleware that validates JWT tokens
func AuthMiddleware(authService *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondWithError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			// Check Bearer prefix
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondWithError(w, http.StatusUnauthorized, "invalid authorization header format")
				return
			}

			token := parts[1]

			// Validate token
			userID, err := authService.ValidateToken(token)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			// Add user ID to request context
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext retrieves the user ID from the request context
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(userIDKey).(int)
	return userID, ok
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	response := models.ErrorResponse{
		Error: message,
	}
	// Simple JSON encoding without importing encoding/json in middleware
	w.Write([]byte(`{"error":"` + response.Error + `"}`))
}
