// oilan/internal/infrastructure/server/server.go
package server

import (
	"net/http"
	"oilan/internal/infrastructure/handlers"
	"time"
)

// NewServer now accepts both API and Page handlers.
func NewServer(api *handlers.APIHandlers, pages *handlers.PageHandlers) *http.Server {
	router := http.NewServeMux()

	// CORRECTED CALL: We call the standalone function from the handlers package,
	// passing it all the necessary dependencies.
	handlers.RegisterRoutes(router, api, pages)

	return &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}