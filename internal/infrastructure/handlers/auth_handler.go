// oilan/internal/infrastructure/handlers/auth_handler.go
package handlers

import (
	"context"
	"fmt"
	"net/http"
	"oilan/internal/domain"
	"time"

	"github.com/markbates/goth/gothic"
)

// BeginAuthHandler initiates the authentication process.
// It redirects the user to the provider's login page (e.g., Google).
func (h *APIHandlers) BeginAuthHandler(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

// AuthCallbackHandler handles the callback from the provider after the user has authenticated.
func (h *APIHandlers) AuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	gothUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintf(w, "Error completing user auth: %v", err)
		return
	}

	// At this point, we have the user's data from Google.
	// Now, we need to save this user to our own database.
	user := &domain.User{
		Provider:   gothUser.Provider,
		ProviderID: gothUser.UserID,
		Email:      gothUser.Email,
		CreatedAt:  time.Now(),
	}

	// Check if this user already exists.
	existingUser, err := h.userRepo.FindByProviderID(context.Background(), user.Provider, user.ProviderID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if existingUser == nil {
		// If user does not exist, save them.
		if err := h.userRepo.Save(context.Background(), user); err != nil {
			http.Error(w, "Failed to save user", http.StatusInternalServerError)
			return
		}
	} else {
		// If user exists, we'll just use their existing ID.
		user.ID = existingUser.ID
	}
	
	// IMPORTANT: Here we will create a session for the user (e.g., JWT token).
	// This will be our next step. For now, we just show the user's info.
	fmt.Fprintf(w, "Login successful! Welcome, %s! Your ID in our system is %d", user.Email, user.ID)
}