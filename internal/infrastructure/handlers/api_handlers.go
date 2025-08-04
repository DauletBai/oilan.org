// oilan/internal/infrastructure/handlers/api_handlers.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"oilan/internal/app/services"
	"oilan/internal/domain/repository"
	"oilan/internal/infrastructure/middleware"
	"strconv"
)

// APIHandlers holds all dependencies for API handlers.
type APIHandlers struct {
	chatService *services.ChatService
	userRepo    repository.UserRepository
}

// NewAPIHandlers creates a new instance of APIHandlers.
func NewAPIHandlers(cs *services.ChatService, ur repository.UserRepository) *APIHandlers {
	return &APIHandlers{
		chatService: cs,
		userRepo:    ur,
	}
}

// --- Helper Functions ---

// writeJSON is a helper for sending JSON responses.
func (h *APIHandlers) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// writeError is a helper for sending JSON error responses.
func (h *APIHandlers) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}

// --- API Handlers ---

// CreateDialogHandler handles requests to create a new dialog.
func (h *APIHandlers) CreateDialogHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user ID that the middleware has placed in the context.
	currentUserID, ok := r.Context().Value(middleware.UserIDContextKey).(int64)
	if !ok {
		h.writeError(w, http.StatusUnauthorized, "Could not identify user")
		return
	}

	var requestBody struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	dialog, err := h.chatService.StartNewDialog(r.Context(), currentUserID, requestBody.Title)
	if err != nil {
		log.Printf("Error creating dialog: %v", err)
		h.writeError(w, http.StatusInternalServerError, "Could not create dialog")
		return
	}

	h.writeJSON(w, http.StatusCreated, dialog)
}

// PostMessageHandler now gets the user ID from the request context.
func (h *APIHandlers) PostMessageHandler(w http.ResponseWriter, r *http.Request) {
	currentUserID, ok := r.Context().Value(middleware.UserIDContextKey).(int64)
	if !ok {
		h.writeError(w, http.StatusUnauthorized, "Could not identify user")
		return
	}

	// Extract dialogID from the URL path, e.g., /api/v1/dialogs/123/messages
	dialogIDStr := r.PathValue("dialogID")
	dialogID, err := strconv.ParseInt(dialogIDStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid dialog ID")
		return
	}

	var requestBody struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if requestBody.Content == "" {
		h.writeError(w, http.StatusBadRequest, "Content cannot be empty")
		return
	}

	aiMessage, err := h.chatService.PostMessage(r.Context(), dialogID, currentUserID, requestBody.Content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeError(w, http.StatusNotFound, "Dialog not found")
			return
		}
		log.Printf("Error posting message: %v", err)
		h.writeError(w, http.StatusInternalServerError, "Failed to process message")
		return
	}

	h.writeJSON(w, http.StatusOK, aiMessage)
}