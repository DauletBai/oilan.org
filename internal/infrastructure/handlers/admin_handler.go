// github.com/DauletBai/oilan.org/internal/infrastructure/handlers/admin_handler.go
package handlers

import (
	"log"
	"net/http"
	"github.com/DauletBai/oilan.org/internal/domain/repository" 
	"github.com/DauletBai/oilan.org/internal/view"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// AdminHandlers holds dependencies for admin page handlers.
type AdminHandlers struct {
	DashboardTemplate  *view.Template
	UsersTemplate      *view.Template 
	DialogsTemplate    *view.Template
	DialogViewTemplate *view.Template
	UserRepo           repository.UserRepository
	DialogRepo         repository.DialogRepository
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

// DialogViewHandler renders a single dialog with all its messages.
func (h *AdminHandlers) DialogViewHandler(w http.ResponseWriter, r *http.Request) {
	// Use chi to get the URL parameter
	dialogIDStr := chi.URLParam(r, "dialogID")
	dialogID, err := strconv.ParseInt(dialogIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid dialog ID", http.StatusBadRequest)
		return
	}

	// Fetch the full dialog with messages using the method we already have
	dialog, err := h.DialogRepo.FindByID(r.Context(), dialogID)
	if err != nil {
		log.Printf("Error finding dialog by ID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if dialog == nil {
		http.Error(w, "Dialog not found", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"title":  "View Dialog",
		"dialog": dialog,
	}

	err = h.DialogViewTemplate.Render(w, "base.html", data)
	if err != nil {
		log.Printf("Error rendering dialog view template: %v", err)
	}
}