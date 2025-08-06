// oilan/internal/infrastructure/middleware/admin.go
package middleware

import (
	"context"
	"net/http"
	"oilan/internal/domain/repository"
)

// AdminMiddleware is a factory that returns a new middleware handler.
// It takes the user repository as a dependency.
func AdminMiddleware(userRepo repository.UserRepository) func(http.Handler) http.Handler {
	// This is the actual middleware that will be returned and used by the router.
	return func(next http.Handler) http.Handler {
		// This is the handler function that runs on every request.
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value(UserIDContextKey).(int64)
			if !ok {
				http.Error(w, "Unauthorized: Not logged in", http.StatusUnauthorized)
				return
			}

			user, err := userRepo.FindByID(context.Background(), userID)
			if err != nil || user == nil {
				http.Error(w, "Unauthorized: User not found", http.StatusUnauthorized)
				return
			}

			if user.Role != "admin" {
				http.Error(w, "Forbidden: You are not an admin", http.StatusForbidden)
				return
			}

			// If all checks pass, proceed to the next handler in the chain.
			next.ServeHTTP(w, r)
		})
	}
}