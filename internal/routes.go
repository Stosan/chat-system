package internal

import (
	"chatsystem/internal/handlers"
	app_midd "chatsystem/internal/middleware"
	ws "chatsystem/internal/websocket"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupWebSocketRoutes(e *echo.Echo, db *gorm.DB, rdb *redis.Client) {
	hub := ws.NewHub()
	wsHandler := handlers.NewWebSocketHandler(hub)
	// Start goroutines to process channels
	go hub.ProcessChatMessages()
	go hub.ProcessPersistMessages()
	e.GET("/ws/server", wsHandler.HandleWebSocket)
}

func ApiRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	e.Use(app_midd.Recover)
	// e.Use(app_midd.SetHeaders)

	chatGroup := e.Group("v1/chat")

	// Initialize handlers
	chatHandler := handlers.NewWebSocketChatHandler(db, rdb)

	// Define routes
	chatGroup.POST("/register", chatHandler.RegisterHandler)
	chatGroup.POST("/send", chatHandler.SendHandler)
	chatGroup.GET("/listen/:userID", chatHandler.ListenHandler)
	chatGroup.DELETE("/disconnect/:userID", chatHandler.DisconnectHandler)
}
