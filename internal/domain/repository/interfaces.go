// github.com/DauletBai/oilan.org/internal/domain/repository/interfaces.go
package repository

import (
	"context"
	"github.com/DauletBai/oilan.org/internal/domain"
)

// UserRepository defines the interface for user data storage.
type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id int64) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByProviderID(ctx context.Context, provider string, providerID string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	GetAll(ctx context.Context) ([]*domain.User, error) 
}

// DialogRepository defines the interface for dialog data storage.
type DialogRepository interface {
	Save(ctx context.Context, dialog *domain.Dialog) error
	FindByID(ctx context.Context, id int64) (*domain.Dialog, error)
	FindAllByUserID(ctx context.Context, userID int64) ([]*domain.Dialog, error)
	GetAll(ctx context.Context) ([]*domain.Dialog, error) 
	AddMessage(ctx context.Context, message *domain.Message) error
}