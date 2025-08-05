// oilan/internal/infrastructure/handlers/pages.go
package handlers

import (
	"log"
	"net/http"
	"oilan/internal/view"
)

// PageHandlers holds dependencies for page rendering handlers.
type PageHandlers struct {
	WelcomeTemplate *view.Template
	ChatTemplate    *view.Template // <-- Add the chat template
}

// WelcomeHandler renders the main welcome page.
func (h *PageHandlers) WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{"title": "Welcome"}
	err := h.WelcomeTemplate.Render(w, "base.html", data)
	if err != nil {
		log.Printf("Error rendering welcome template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ChatHandler renders the chat page.
func (h *PageHandlers) ChatHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{"title": "Chat"}
	err := h.ChatTemplate.Render(w, "base.html", data)
	if err != nil {
		log.Printf("Error rendering chat template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}