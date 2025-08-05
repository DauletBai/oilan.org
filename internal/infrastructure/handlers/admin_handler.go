// oilan/internal/infrastructure/handlers/admin_handler.go
package handlers

import (
	"log"
	"net/http"
	"oilan/internal/view" 
)

// AdminHandlers holds dependencies for admin page handlers.
type AdminHandlers struct {
	DashboardTemplate *view.Template
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