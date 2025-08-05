// oilan/internal/infrastructure/server/server.go
package server

import (
	"net/http"
	"oilan/internal/domain/repository"
	"oilan/internal/infrastructure/handlers"
	"time"
)

// NewServer now accepts all handler types and the user repository for middleware.
func NewServer(api *handlers.APIHandlers, pages *handlers.PageHandlers, admin *handlers.AdminHandlers, userRepo repository.UserRepository) *http.Server {
	router := http.NewServeMux()

	// Pass all dependencies to the router function.
	handlers.RegisterRoutes(router, api, pages, admin, userRepo)

	return &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}