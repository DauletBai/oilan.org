// github.com/DauletBai/oilan.org/internal/infrastructure/handlers/api_handlers.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"github.com/DauletBai/oilan.org/internal/app/services"
	"github.com/DauletBai/oilan.org/internal/domain/repository"
	"github.com/DauletBai/oilan.org/internal/infrastructure/middleware"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// APIHandlers holds all dependencies for API handlers.
type APIHandlers struct {
	chatService *services.ChatService
	userRepo    repository.UserRepository
	dialogRepo  repository.DialogRepository
}

// NewAPIHandlers creates a new instance of APIHandlers.
func NewAPIHandlers(cs *services.ChatService, ur repository.UserRepository, dr repository.DialogRepository) *APIHandlers {
	return &APIHandlers{
		chatService: cs,
		userRepo:    ur,
		dialogRepo:  dr,
	}
}

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

// GetSessionInfoHandler provides the frontend with essential session data.
func (h *APIHandlers) GetSessionInfoHandler(w http.ResponseWriter, r *http.Request) {
	// The AuthMiddleware has already run and placed the user ID in the context.
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(int64)
	if !ok {
		// This should theoretically never happen if the middleware is applied correctly.
		h.writeError(w, http.StatusUnauthorized, "No valid session found")
		return
	}

	// We can add more user info here if needed in the future.
	sessionData := map[string]interface{}{
		"userID": userID,
	}

	h.writeJSON(w, http.StatusOK, sessionData)
}

// GetDialogsHandler returns a list of all dialogs for the authenticated user.
func (h *APIHandlers) GetDialogsHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.UserIDContextKey).(int64)
	dialogs, err := h.dialogRepo.FindAllByUserID(r.Context(), userID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Could not retrieve dialogs")
		return
	}
	h.writeJSON(w, http.StatusOK, dialogs)
}

// CreateDialogHandler handles requests to create a new dialog.
func (h *APIHandlers) CreateDialogHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.UserIDContextKey).(int64)

	// Get the user ID that the middleware has placed in the context.
	// currentUserID, ok := r.Context().Value(middleware.UserIDContextKey).(int64)
	// if !ok {
	// 	h.writeError(w, http.StatusUnauthorized, "Could not identify user")
	// 	return
	// }

	var requestBody struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	dialog, err := h.chatService.StartNewDialog(r.Context(), userID, requestBody.Title)
	if err != nil {
		log.Printf("Error creating dialog: %v", err)
		h.writeError(w, http.StatusInternalServerError, "Could not create dialog")
		return
	}

	h.writeJSON(w, http.StatusCreated, dialog)
}

// PostMessageHandler now gets the user ID from the request context.
func (h *APIHandlers) PostMessageHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.UserIDContextKey).(int64)

	// currentUserID, ok := r.Context().Value(middleware.UserIDContextKey).(int64)
	// if !ok {
	// 	h.writeError(w, http.StatusUnauthorized, "Could not identify user")
	// 	return
	// }

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

	aiMessage, err := h.chatService.PostMessage(r.Context(), dialogID, userID, requestBody.Content)
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

// GetDialogByIDHandler returns a single dialog with all its messages.
func (h *APIHandlers) GetDialogByIDHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.UserIDContextKey).(int64)

	dialogIDStr := chi.URLParam(r, "dialogID") // Using chi to get URL parameter
	dialogID, err := strconv.ParseInt(dialogIDStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid dialog ID")
		return
	}

	dialog, err := h.dialogRepo.FindByID(r.Context(), dialogID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Could not retrieve dialog")
		return
	}
	if dialog == nil {
		h.writeError(w, http.StatusNotFound, "Dialog not found")
		return
	}

	// Security check: ensure the user owns this dialog.
	if dialog.UserID != userID {
		h.writeError(w, http.StatusForbidden, "Access denied")
		return
	}

	h.writeJSON(w, http.StatusOK, dialog)
}
