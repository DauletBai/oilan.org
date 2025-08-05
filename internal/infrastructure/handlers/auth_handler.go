// oilan/internal/infrastructure/handlers/auth_handler.go
package handlers

import (
	"context"
	//"fmt"
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

// AuthCallbackHandler now sets a secure cookie and redirects to the chat page.
func (h *APIHandlers) AuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	gothUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		// In case of error, redirect to the home page with an error message
		http.Redirect(w, r, "/?error=auth_failed", http.StatusTemporaryRedirect)
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
	
	tokenString, err := auth.GenerateToken(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Set the token in a secure, HttpOnly cookie.
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    tokenString,
		Expires:  time.Now().Add(72 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil, // Use secure cookies in production (HTTPS)
		SameSite: http.SameSiteLaxMode,
	})

	// Redirect the user to the chat page.
	http.Redirect(w, r, "/chat", http.StatusTemporaryRedirect)
}
