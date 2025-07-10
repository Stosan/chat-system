package handlers

import (
	"log"
	"net/http"
	"time"

	"chatsystem/internal/models"
	ws "chatsystem/internal/websocket"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WebSocketHandler struct {
	hub *ws.Hub
}

func NewWebSocketHandler(hub *ws.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
	}
}

// HandleWebSocket handles WebSocket connections with Echo
func (h *WebSocketHandler) HandleWebSocket(c echo.Context) error {
	// Upgrade HTTP connection to WebSocket
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Add your origin validation logic here
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return err
	}
	defer conn.Close()

	// Handle the WebSocket connection
	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		msg.Timestamp = time.Now()

		switch msg.Type {
		case "new_client":
			if h.hub.RegisterUser(msg.Sender, conn) {
				response := models.Message{
					Text:      "Registration successful",
					Sender:    "server",
					Receiver:  msg.Sender,
					Type:      "registration_success",
					Timestamp: time.Now(),
				}
				conn.WriteJSON(response)

				// Fetch historical messages in goroutine
				// go h.hub.FetchHistoricalMessages(msg.Sender)
			} else {
				response := models.Message{
					Text:      "User already exists",
					Sender:    "server",
					Receiver:  msg.Sender,
					Type:      "registration_error",
					Timestamp: time.Now(),
				}
				conn.WriteJSON(response)
			}

		case "chat":
			// Send to chat channel for delivery
			h.hub.SendToChat(msg)
			// Send to persist channel for storage
			h.hub.SendToPersist(msg)

		case "session_end":
			h.hub.RemoveUser(msg.Sender)
			return nil
		}
	}

	return nil
}
