// oilan/internal/infrastructure/handlers/websocket_handler.go
package handlers

import (
	"log"
	"net/http"
	"oilan/internal/infrastructure/middleware" 
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// In production, you should check the origin of the request.
	// For development, we can allow any origin.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ServeWs handles websocket requests from the client.
func (h *APIHandlers) ServeWs(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Get user ID from the context (thanks to our middleware!)
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(int64)
	if !ok {
		log.Println("Could not get user ID from context")
		return
	}

	// Get dialog ID from the query parameter, e.g., /ws/chat?dialogID=123
	dialogIDStr := r.URL.Query().Get("dialogID")
	dialogID, err := strconv.ParseInt(dialogIDStr, 10, 64)
	if err != nil {
		log.Println("Invalid dialog ID")
		return
	}
	
	log.Printf("User %d connected to dialog %d via WebSocket", userID, dialogID)

	// This is the main loop for the connection.
	for {
		// Read a message from the browser.
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		// Log the message for now.
		log.Printf("Received from user %d: %s", userID, msg)

		// Here, we will call our chatService to get the AI response.
		aiResponse, err := h.chatService.PostMessage(r.Context(), dialogID, userID, string(msg))
		if err != nil {
			log.Println("ChatService error:", err)
			// Send an error message back to the client.
			conn.WriteMessage(websocket.TextMessage, []byte("Sorry, an error occurred."))
			continue
		}

		// Send the AI's response back to the browser.
		if err := conn.WriteMessage(websocket.TextMessage, []byte(aiResponse.Content)); err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}