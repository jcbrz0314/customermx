package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/customermx/backend/internal/config"
	"github.com/customermx/backend/internal/infra/security"
)

type contextKey string

const (
	// UserContextKey is the key for user claims in context
	UserContextKey contextKey = "user"
)

// AuthMiddleware validates JWT tokens and adds user claims to context
func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	jwtService := security.NewJWTService(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.AccessTTL,
		cfg.JWT.RefreshTTL,
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Extract token from "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			token := parts[1]
			claims, err := jwtService.ValidateAccessToken(token)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Check if user is active
			if !claims.IsActive {
				http.Error(w, "User account is inactive", http.StatusForbidden)
				return
			}

			// Add claims to context
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext retrieves user claims from context
func GetUserFromContext(ctx context.Context) (*security.Claims, bool) {
	claims, ok := ctx.Value(UserContextKey).(*security.Claims)
	return claims, ok
}
