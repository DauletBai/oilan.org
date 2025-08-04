// oilan/internal/infrastructure/repository/postgres/postgres.go
package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

// NewConnection creates a new database connection.
func NewConnection() (*sql.DB, error) {
	// Build connection string from environment variables
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify the connection is alive.
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}