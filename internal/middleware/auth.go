package middleware

import (
	"context"
	"log"
	"max-odyssey-backend/internal/models"
	"max-odyssey-backend/internal/service"
	"net/http"
	"strings"
)

type contextKey string

const UserContextKey contextKey = "user"

func Auth(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Println("No Authorization header found")
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Check if the header is in the correct format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				log.Println("Invalid Authorization format:", authHeader)
				http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
				return
			}

			// Validate the token
			user, err := authService.ValidateToken(parts[1])
			if err != nil {
				log.Printf("Token validation error: %v", err)
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			log.Printf("Authenticated user: %s (ID: %d)", user.Username, user.ID)

			// Add the user to the request context
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext gets the user from the request context
func GetUserFromContext(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value(UserContextKey).(*models.User)
	return user, ok
}
