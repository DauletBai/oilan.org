// oilan.org/internal/domain/user.go
package domain

import (
	"time"
)

// User represents a user in the system.
type User struct {
	ID        int64     `json:"id"`         // Unique identifier
	Provider  string    `json:"provider"`   // e.g., "google", "microsoft"
	ProviderID string   `json:"provider_id"`// User ID from the provider
	Email     string    `json:"email"`      // User's email, verified by provider
	Role      string    `json:"role"`       // "user", "admin")
	CreatedAt time.Time `json:"created_at"` // Timestamp of user creation
}