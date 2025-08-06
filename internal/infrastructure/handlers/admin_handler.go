// oilan/internal/infrastructure/handlers/admin_handler.go
package handlers

import (
	"log"
	"net/http"
	"oilan/internal/domain/repository" 
	"oilan/internal/view"
)

// AdminHandlers holds dependencies for admin page handlers.
type AdminHandlers struct {
	DashboardTemplate *view.Template
	UsersTemplate     *view.Template 
	DialogsTemplate   *view.Template
	UserRepo          repository.UserRepository
	DialogRepo        repository.DialogRepository
}

// DashboardHandler renders the main admin dashboard page.
func (h *AdminHandlers) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"title": "Dashboard",
	}
	err := h.DashboardTemplate.Render(w, "base.html", data)
	if err != nil {
		log.Printf("Error rendering dashboard template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// UsersHandler renders the user management page.
func (h *AdminHandlers) UsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.UserRepo.GetAll(r.Context())
	if err != nil {
		log.Printf("Error getting all users: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"title": "User Management",
		"users": users,
	}

	err = h.UsersTemplate.Render(w, "base.html", data)
	if err != nil {
		log.Printf("Error rendering users template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// DialogsHandler renders the dialog management page.
func (h *AdminHandlers) DialogsHandler(w http.ResponseWriter, r *http.Request) {
	dialogs, err := h.DialogRepo.GetAll(r.Context())
	if err != nil {
		log.Printf("Error getting all dialogs: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"title":   "Dialogs Management",
		"dialogs": dialogs,
	}

	err = h.DialogsTemplate.Render(w, "base.html", data)
	if err != nil {
		log.Printf("Error rendering dialogs template: %v", err)
	}
}