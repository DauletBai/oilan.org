// oilan/internal/infrastructure/handlers/routes.go
package handlers

import "net/http"

// RegisterRoutes registers all application routes.
func (h *APIHandlers) RegisterRoutes(mux *http.ServeMux) {
	// --- Static File Server for Frontend ---
	// Create a file server to serve files out of the ./web/static directory.
	// This will automatically serve index.html for the root path.
	fileServer := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/", fileServer)

	// --- Authentication Handlers ---
	// These handlers are for the OAuth flow.
	mux.HandleFunc("/auth/{provider}/callback", h.AuthCallbackHandler)
	mux.HandleFunc("/auth/{provider}", h.BeginAuthHandler)

	// --- Secure API Handlers ---
	// These handlers are for the application's core functionality.
	// mux.Handle("POST /api/v1/dialogs", middleware.AuthMiddleware(http.HandlerFunc(h.CreateDialogHandler)))
	// mux.Handle("POST /api/v1/dialogs/{dialogID}/messages", middleware.AuthMiddleware(http.HandlerFunc(h.PostMessageHandler)))
}