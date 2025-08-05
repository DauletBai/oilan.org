// oilan/internal/infrastructure/handlers/websocket_handler.go
package handlers

import (
	"log"
	"net/http"
	"oilan/internal/infrastructure/middleware" 

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // For development
}

// ServeWs handles websocket requests from the client.
func (h *APIHandlers) ServeWs(w http.ResponseWriter, r *http.Request) {
	// --- THE CRITICAL FIX ---
	// The AuthMiddleware has already run. We now get the userID from the context.
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(int64)
	if !ok {
		// This will only happen if middleware is not applied correctly.
		http.Error(w, "Unauthorized: No user ID in context", http.StatusUnauthorized)
		return
	}

	// We no longer need the dialogID from the query. A new dialog will be created on connection.
	
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// --- Create a new dialog for this WebSocket session ---
	dialog, err := h.chatService.StartNewDialog(r.Context(), userID, "New WebSocket Chat")
	if err != nil {
		log.Printf("Failed to create new dialog for user %d: %v", userID, err)
		return
	}
	dialogID := dialog.ID
	log.Printf("User %d connected to new dialog %d via WebSocket", userID, dialogID)
	
	// Send a welcome message.
	conn.WriteMessage(websocket.TextMessage, []byte("Hello! I am ready. How can I help you today?"))

	// Main loop for the connection.
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("User %d disconnected from dialog %d", userID, dialogID)
			break
		}

		aiResponse, err := h.chatService.PostMessage(r.Context(), dialogID, userID, string(msg))
		if err != nil {
			log.Println("ChatService error:", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Sorry, an error occurred."))
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, []byte(aiResponse.Content)); err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}