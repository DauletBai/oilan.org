// oilan/internal/infrastructure/server/server.go
package server

import (
	"net/http"
	"oilan/internal/domain/repository"
	"oilan/internal/infrastructure/handlers"
	"time"
)

// NewServer now uses the router returned by RegisterRoutes.
func NewServer(api *handlers.APIHandlers, pages *handlers.PageHandlers, admin *handlers.AdminHandlers, userRepo repository.UserRepository) *http.Server {
	// The router is now configured inside RegisterRoutes
	router := handlers.RegisterRoutes(api, pages, admin, userRepo)

	return &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}