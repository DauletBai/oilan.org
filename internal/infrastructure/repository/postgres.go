// oilan/internal/infrastructure/repository/postgres/postgres.go
package postgres 

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" 
	"os"
)

// NewConnection создает новое соединение с базой данных.
func NewConnection() (*sql.DB, error) {
	// Собираем строку подключения из переменных окружения
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

	// Проверяем, что соединение действительно установлено
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}