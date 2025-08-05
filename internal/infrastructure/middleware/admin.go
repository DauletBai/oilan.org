// oilan/internal/infrastructure/middleware/admin.go
package middleware

import (
	"context"
	"net/http"
	"oilan/internal/domain/repository"
)

// AdminMiddleware checks if the authenticated user has the 'admin' role.
// It must run *after* the AuthMiddleware.
func AdminMiddleware(next http.Handler, userRepo repository.UserRepository) http.Handler {
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

		// If all checks pass, proceed to the next handler.
		next.ServeHTTP(w, r)
	})
}