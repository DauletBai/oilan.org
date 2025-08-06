// oilan/internal/infrastructure/middleware/auth.go
// package middleware

// import (
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"oilan/internal/auth"
// 	"os"
// 	"strings"

// 	"github.com/golang-jwt/jwt/v5"
// )

// type contextKey string

// const UserIDContextKey = contextKey("userID")

// // AuthMiddleware verifies the JWT token from the request header or URL query.
// func AuthMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// First, try to get the token from the "Authorization" header.
// 		tokenString := ""
// 		authHeader := r.Header.Get("Authorization")
// 		if authHeader != "" {
// 			headerParts := strings.Split(authHeader, " ")
// 			if len(headerParts) == 2 && headerParts[0] == "Bearer" {
// 				tokenString = headerParts[1]
// 			}
// 		}

// 		// If not in header (e.g., a WebSocket request), try to get it from a query parameter.
// 		if tokenString == "" {
// 			tokenString = r.URL.Query().Get("token")
// 		}

// 		if tokenString == "" {
// 			http.Error(w, "Authorization token required", http.StatusUnauthorized)
// 			return
// 		}

// 		secretKey := []byte(os.Getenv("SESSION_SECRET"))
// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 			}
// 			return secretKey, nil
// 		})

// 		if err != nil || !token.Valid {
// 			http.Error(w, "Invalid token", http.StatusUnauthorized)
// 			return
// 		}

// 		claims, ok := token.Claims.(jwt.MapClaims)
// 		if !ok {
// 			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
// 			return
// 		}

// 		userIDFloat, ok := claims["sub"].(float64)
// 		if !ok {
// 			http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
// 			return
// 		}
// 		userID := int64(userIDFloat)

// 		ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

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