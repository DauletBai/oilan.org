// oilan/internal/infrastructure/handlers/routes.go
package handlers

import "net/http"

// RegisterRoutes registers all application routes.
func (h *APIHandlers) RegisterRoutes(mux *http.ServeMux) {
	// Page handlers
	mux.HandleFunc("/", HomeHandler)

	// Authentication handlers
	mux.HandleFunc("/auth/{provider}/callback", h.AuthCallbackHandler)
	mux.HandleFunc("/auth/{provider}", h.BeginAuthHandler)

	// API handlers
	mux.HandleFunc("POST /api/v1/dialogs", h.CreateDialogHandler)
	mux.HandleFunc("POST /api/v1/dialogs/{dialogID}/messages", h.PostMessageHandler)
}