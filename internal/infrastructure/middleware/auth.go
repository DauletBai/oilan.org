// oilan/internal/infrastructure/middleware/auth.go
package middleware

import (
	"context"
	"fmt"
	"net/http"
	//"oilan/internal/auth" 
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDContextKey = contextKey("userID")

// AuthMiddleware now verifies the JWT token from the secure HttpOnly cookie.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Get the cookie from the request.
		cookie, err := r.Cookie("jwt_token")
		if err != nil {
			// If no cookie is found, the user is not authenticated.
			// We redirect them to the login page instead of showing an error.
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		tokenString := cookie.Value
		secretKey := []byte(os.Getenv("SESSION_SECRET"))

		// 2. Parse and validate the token.
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		// 3. Check for errors or if the token is invalid.
		if err != nil || !token.Valid {
			// If token is invalid, delete the bad cookie and redirect to login.
			http.SetCookie(w, &http.Cookie{Name: "jwt_token", Value: "", Path: "/", MaxAge: -1})
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		
		// 4. Extract user ID from the token's claims.
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok { // Should not happen with our own tokens
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		userIDFloat, ok := claims["sub"].(float64)
		if !ok { // Should not happen with our own tokens
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		userID := int64(userIDFloat)

		// 5. Add the user ID to the request's context.
		ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
		
		// 6. Call the next handler in the chain with the updated context.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}