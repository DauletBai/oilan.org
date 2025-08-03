// oilan.org/internal/domain/dialog.go
package domain

import "time"

// Role defines who is the author of a message.
type Role string

const (
	RoleUser Role = "user"
	RoleAI   Role = "ai"
)

// Message represents a single message within a dialog.
type Message struct {
	ID        int64     `json:"id"`
	DialogID  int64     `json:"dialog_id"`
	Role      Role      `json:"role"`      // "user" or "ai"
	Content   string    `json:"content"`   // The text of the message
	CreatedAt time.Time `json:"created_at"`
}

// Dialog represents a complete conversation session for a user.
type Dialog struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`     // A title for the dialog, can be auto-generated
	Messages  []Message `json:"messages"`  // The list of all messages in this dialog
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}