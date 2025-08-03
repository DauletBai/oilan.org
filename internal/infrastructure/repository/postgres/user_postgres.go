// oilan/internal/infrastructure/repository/postgres/user_postgres.go
package postgres

import (
	"context"
	"database/sql"
	"errors" 
	"oilan/internal/domain"
	"oilan/internal/domain/repository"
)

// userRepo implements the repository.UserRepository interface.
type userRepo struct {
	db *sql.DB
}

// NewUserRepository creates a new instance of the user repository.
func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepo{db: db}
}

// Save saves a new user or updates an existing one based on provider and provider_id.
func (r *userRepo) Save(ctx context.Context, user *domain.User) error {
	query := `
        INSERT INTO users (provider, provider_id, email, created_at)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (provider, provider_id) DO UPDATE SET email = EXCLUDED.email
        RETURNING id;
    `
	return r.db.QueryRowContext(ctx, query, user.Provider, user.ProviderID, user.Email, user.CreatedAt).Scan(&user.ID)
}

// FindByID finds a user by their unique internal ID.
func (r *userRepo) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT id, provider, provider_id, email, created_at FROM users WHERE id = $1;`
	
	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Provider,
		&user.ProviderID,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// This is a standard way to handle "not found" cases.
			return nil, nil 
		}
		return nil, err
	}

	return user, nil
}

// FindByProviderID finds a user by their provider and provider-specific ID.
func (r *userRepo) FindByProviderID(ctx context.Context, provider string, providerID string) (*domain.User, error) {
	query := `SELECT id, provider, provider_id, email, created_at FROM users WHERE provider = $1 AND provider_id = $2;`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, provider, providerID).Scan(
		&user.ID,
		&user.Provider,
		&user.ProviderID,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found, which is not a system error.
		}
		return nil, err
	}

	return user, nil
}