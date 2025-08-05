// oilan/internal/infrastructure/handlers/routes.go
package handlers

import (
	"net/http"
	"oilan/internal/domain/repository"
	"oilan/internal/infrastructure/middleware"
)

// RegisterRoutes is a standalone function that registers all application routes.
func RegisterRoutes(mux *http.ServeMux, api *APIHandlers, pages *PageHandlers, admin *AdminHandlers, userRepo repository.UserRepository) {
	// --- Static File Server ---
	fileServer := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// --- Page Handlers ---
	mux.HandleFunc("/", pages.WelcomeHandler)
	mux.HandleFunc("/chat", pages.ChatHandler)

	// --- Authentication Handlers ---
	mux.HandleFunc("/auth/{provider}/callback", api.AuthCallbackHandler)
	mux.HandleFunc("/auth/{provider}", api.BeginAuthHandler)
	
	// --- Secure API Handlers ---
	// Create a sub-router for our authenticated API v1
	apiV1 := http.NewServeMux()
	apiV1.HandleFunc("POST /dialogs", api.CreateDialogHandler)
	apiV1.HandleFunc("POST /dialogs/{dialogID}/messages", api.PostMessageHandler)
	apiV1.HandleFunc("GET /session", api.GetSessionInfoHandler)

	// Apply the AuthMiddleware to all /api/v1 routes
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", middleware.AuthMiddleware(apiV1)))

	// --- WebSocket Handler ---
	mux.Handle("/ws/chat", middleware.AuthMiddleware(http.HandlerFunc(api.ServeWs)))

	// --- Admin Routes ---
	adminRouter := http.NewServeMux()
	adminRouter.HandleFunc("/admin/dashboard", admin.DashboardHandler)
	
	// Protect all admin routes with two layers of security
	mux.Handle("/admin/", http.StripPrefix("/admin", 
		middleware.AuthMiddleware(
			middleware.AdminMiddleware(adminRouter, userRepo),
		),
	))
}