// oilan/internal/infrastructure/server/server.go
package server

import (
	"net/http"
	"time"
	"oilan/internal/infrastructure/handlers"
)

// NewServer создает и настраивает новый HTTP сервер.
func NewServer() *http.Server {
    // Создаем новый роутер (мультиплексор)
	router := http.NewServeMux()

    // Регистрируем все наши маршруты (URL)
	handlers.RegisterRoutes(router)

	return &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}