// oilan/internal/infrastructure/handlers/routes.go
package handlers

import (
	"net/http"
	"oilan/internal/infrastructure/middleware"
)

// RegisterRoutes is a standalone function that registers all application routes.
func RegisterRoutes(mux *http.ServeMux, api *APIHandlers, pages *PageHandlers) {
	// --- Static File Server ---
	fileServer := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// --- Page Handlers ---
	mux.HandleFunc("/", pages.WelcomeHandler)

	// --- Authentication Handlers ---
	mux.HandleFunc("/auth/{provider}/callback", api.AuthCallbackHandler)
	mux.HandleFunc("/auth/{provider}", api.BeginAuthHandler)
	
    // --- WebSocket Handler ---
	mux.Handle("/ws/chat", middleware.AuthMiddleware(http.HandlerFunc(api.ServeWs)))
}