// oilan/internal/infrastructure/server/server.go
package server

import (
	"net/http"
	"time"
	"oilan/internal/infrastructure/handlers"
)

// NewServer creates and configures a new HTTP server.
func NewServer(api *handlers.APIHandlers) *http.Server {
	router := http.NewServeMux()

	// CORRECTED LINE: We now call RegisterRoutes as a method on the api object.
	api.RegisterRoutes(router)

	return &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}