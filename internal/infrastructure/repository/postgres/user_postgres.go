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

// FindByEmail finds a user by their email address.
func (r *userRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, provider, provider_id, email, role, created_at FROM users WHERE email = $1;`
	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Provider, &user.ProviderID, &user.Email, &user.Role, &user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found is a valid case
		}
		return nil, err
	}
	return user, nil
}

// Update updates an existing user's data (e.g., their role).
func (r *userRepo) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET role = $1 WHERE id = $2;`
	_, err := r.db.ExecContext(ctx, query, user.Role, user.ID)
	return err
}

// Save saves a new user or updates an existing one based on provider and provider_id.
func (r *userRepo) Save(ctx context.Context, user *domain.User) error {
	query := `
        INSERT INTO users (provider, provider_id, email, role, created_at)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (provider, provider_id) DO UPDATE SET email = EXCLUDED.email
        RETURNING id;
    `
	return r.db.QueryRowContext(ctx, query, user.Provider, user.ProviderID, user.Email, user.Role, user.CreatedAt).Scan(&user.ID)
}

// FindByID finds a user by their unique internal ID.
func (r *userRepo) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT id, provider, provider_id, email, role, created_at FROM users WHERE id = $1;`
	
	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Provider,
		&user.ProviderID,
		&user.Email,
		&user.Role,
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
	query := `SELECT id, provider, provider_id, email, role, created_at FROM users WHERE provider = $1 AND provider_id = $2;`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, provider, providerID).Scan(
		&user.ID,
		&user.Provider,
		&user.ProviderID,
		&user.Email,
		&user.Role,
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

// GetAll retrieves all users from the database.
func (r *userRepo) GetAll(ctx context.Context) ([]*domain.User, error) {
	query := `SELECT id, provider, provider_id, email, role, created_at FROM users ORDER BY created_at DESC;`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(
			&user.ID, &user.Provider, &user.ProviderID, &user.Email, &user.Role, &user.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}