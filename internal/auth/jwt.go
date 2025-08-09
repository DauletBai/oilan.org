// github.com/DauletBai/oilan.org/internal/auth/jwt.go
package auth

import (
	"github.com/DauletBai/oilan.org/internal/domain"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a new JWT token for a given user.
func GenerateToken(user *domain.User) (string, error) {
	// Create the claims for the token
	claims := jwt.MapClaims{
		"sub":   user.ID, // Subject (who the token is for)
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(), // Expiration time (3 days)
		"iat":   time.Now().Unix(),                      // Issued at
	}

	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret key
	// This secret key MUST be the same as the one in your docker-compose.yml
	secretKey := []byte(os.Getenv("SESSION_SECRET"))
	return token.SignedString(secretKey)
}