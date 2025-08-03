// oilan/cmd/server/main.go
package main

import (
	"log"
	"net/http"
	"oilan/internal/infrastructure/repository"
	"oilan/internal/infrastructure/server"
)

func main() {
    // Инициализируем соединение с базой данных
	db, err := postgres.NewConnection()
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer db.Close() // Гарантируем, что соединение будет закрыто при завершении работы

	log.Println("Database connection successful")

    // В будущем мы передадим 'db' в наши сервисы и репозитории
	// userRepo := postgres.NewUserRepository(db)
	// chatService := services.NewChatService(userRepo, ...)

	srv := server.NewServer()
	
	log.Println("Starting server on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("could not listen on :8080: %v\n", err)
	}
}