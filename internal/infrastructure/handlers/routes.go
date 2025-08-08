// oilan/internal/infrastructure/handlers/routes.go
package handlers

import (
	"net/http"
	"oilan/internal/domain/repository"
	"oilan/internal/infrastructure/middleware"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes now uses a cleaner structure for middleware.
func RegisterRoutes(api *APIHandlers, pages *PageHandlers, admin *AdminHandlers, userRepo repository.UserRepository) http.Handler {
	r := chi.NewRouter()

	// --- Public Routes ---
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))
	r.Get("/", pages.WelcomeHandler)
	r.Get("/auth/{provider}", api.BeginAuthHandler)
	r.Get("/auth/{provider}/callback", api.AuthCallbackHandler)

	// --- Authenticated Routes ---
	// All routes inside this group will first pass through AuthMiddleware.
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		// Regular authenticated pages
		r.Get("/chat", pages.ChatHandler)
		r.Get("/ws/chat", api.ServeWs)
		
		// Authenticated API endpoints
		r.Route("/api/v1", func(r chi.Router) {
			r.Get("/session", api.GetSessionInfoHandler)
			r.Post("/dialogs", api.CreateDialogHandler)
			r.Get("/dialogs", api.GetDialogsHandler)
			r.Post("/dialogs/{dialogID}/messages", api.PostMessageHandler)
			r.Get("/dialogs/{dialogID}", api.GetDialogByIDHandler)
		})

		// --- Admin Routes ---
		// This sub-group has an additional AdminMiddleware.
		r.Route("/admin", func(r chi.Router) {
			r.Use(middleware.AdminMiddleware(userRepo))
			r.Get("/dashboard", admin.DashboardHandler)
			r.Get("/users", admin.UsersHandler)
			r.Get("/dialogs", admin.DialogsHandler)
		})
	})

	return r
}