// oilan/internal/infrastructure/handlers/websocket_handler.go
package handlers

import (
	"log"
	"net/http"
	"oilan/internal/infrastructure/middleware"
	//"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (h *APIHandlers) ServeWs(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(int64)
	if !ok {
		http.Error(w, "Unauthorized: No user ID in context", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	dialog, err := h.chatService.StartNewDialog(r.Context(), userID, "New WebSocket Chat")
	if err != nil {
		log.Printf("Failed to create new dialog for user %d: %v", userID, err)
		return
	}
	dialogID := dialog.ID
	log.Printf("User %d connected to new dialog %d via WebSocket", userID, dialogID)
	
	conn.WriteMessage(websocket.TextMessage, []byte("Hello! I am ready. How can I help you today?"))

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("User %d disconnected from dialog %d", userID, dialogID)
			break
		}

		// --- THE CRITICAL FIX ---
		// We now call PostMessage, which is designed to handle the full cycle.
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