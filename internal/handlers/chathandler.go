package handlers

import (
	"chatsystem/internal/services"
	ws "chatsystem/internal/websocket"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type WebSocketChatHandler struct {
	clientManager *ws.ClientManager
	chatService   *services.ChatService
}

func NewWebSocketChatHandler(db *gorm.DB, rdb *redis.Client) *WebSocketChatHandler {
	return &WebSocketChatHandler{
		clientManager: ws.NewClientManager(),
		chatService:   services.NewChatService(db),
	}
}

func (h WebSocketChatHandler) RegisterHandler(c echo.Context) error {
	type Request struct {
		UserID string `json:"user_id" validate:"required"`
	}

	req := new(Request)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if req.UserID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user_id is required")
	}

	// Get or create client
	client := h.clientManager.GetOrCreateClient(req.UserID)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to WebSocket server
	if err := client.ChatConnect(ctx); err != nil {
		h.clientManager.RemoveClient(req.UserID)
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprintf("Failed to connect: %v", err))
	}

	// Register with server
	if err := client.ChatRegister(); err != nil {
		h.clientManager.RemoveClient(req.UserID)
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprintf("Registration failed: %v", err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "registered",
		"user_id": req.UserID,
	})
}

func (h WebSocketChatHandler) SendHandler(c echo.Context) error {
	type Request struct {
		UserID   string `json:"user_id" validate:"required"`
		Receiver string `json:"receiver" validate:"required"`
		Text     string `json:"text" validate:"required"`
	}

	req := new(Request)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if req.UserID == "" || req.Receiver == "" || req.Text == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user_id, receiver, and text are required")
	}

	client, exists := h.clientManager.GetClient(req.UserID)
	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, "User not registered")
	}

	if err := client.SendMessage(req.Receiver, req.Text); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprintf("Failed to send message: %v", err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":    "sent",
		"timestamp": time.Now(),
	})
}

func (h WebSocketChatHandler) ListenHandler(c echo.Context) error {
	userID := c.Param("userID")
	if userID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "userID parameter is required")
	}

	client, exists := h.clientManager.GetClient(userID)
	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, "User not registered")
	}

	// Upgrade to WebSocket
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// In production, implement proper origin checking
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return fmt.Errorf("failed to upgrade to WebSocket: %w", err)
	}
	defer ws.Close()

	// Set up ping/pong for connection health
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Send ping periodically
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Message relay loop
	for {
		select {
		case msg := <-client.GetMessages():
			// Only send messages intended for this user
			if msg.Receiver == userID || msg.Type == "broadcast" {
				ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err := ws.WriteJSON(msg); err != nil {
					log.Printf("Failed to send message to %s: %v", userID, err)
					return nil
				}
			}

		case <-ticker.C:
			ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Failed to send ping to %s: %v", userID, err)
				return nil
			}
		}
	}
}

func (h WebSocketChatHandler) DisconnectHandler(c echo.Context) error {
	userID := c.Param("userID")
	if userID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "userID parameter is required")
	}

	h.clientManager.RemoveClient(userID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "disconnected",
		"user_id": userID,
	})
}
