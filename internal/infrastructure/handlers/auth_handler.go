// oilan/internal/infrastructure/handlers/auth_handler.go
package handlers

import (
	"context"
	"fmt"
	"net/http"
	"oilan/internal/auth" 
	"oilan/internal/domain"
	"time"

	"github.com/markbates/goth/gothic"
)

// ... (BeginAuthHandler remains the same) ...
func (h *APIHandlers) BeginAuthHandler(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

// AuthCallbackHandler now generates a JWT token upon successful login.
func (h *APIHandlers) AuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	gothUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintf(w, "Error completing user auth: %v", err)
		return
	}

	user := &domain.User{
		Provider:   gothUser.Provider,
		ProviderID: gothUser.UserID,
		Email:      gothUser.Email,
	}

	existingUser, err := h.userRepo.FindByProviderID(context.Background(), user.Provider, user.ProviderID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Database error")
		return
	}

	if existingUser == nil {
		user.CreatedAt = time.Now()
		if err := h.userRepo.Save(context.Background(), user); err != nil {
			h.writeError(w, http.StatusInternalServerError, "Failed to save user")
			return
		}
	} else {
		user = existingUser
	}
	
	// Generate JWT Token
	tokenString, err := auth.GenerateToken(user)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// For an actual application, you would redirect the user to the frontend
	// with the token, e.g., http.Redirect(w, r, "http://localhost:3000?token="+tokenString, http.StatusTemporaryRedirect)
	// For now, we will just display it.
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"token": "%s"}`, tokenString)))
}