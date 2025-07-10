package internal

import (
	"chatsystem/internal/config"
	"chatsystem/internal/middleware"
	"chatsystem/pkg/database"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	e_mid "github.com/labstack/echo/v4/middleware"
)

func Start() *echo.Echo {
	log.Println("🟢🔧 Starting Risigner Chat Server")
	// Initialize database connection
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	redisdb, err := database.ConnectRedis()
	if err != nil {
		log.Fatalf("Failed to connect to redis database: %v", err)
	}
	e := echo.New()
	//CORS & Middleware
	e.Use(middleware.CORSMiddleware())
	e.Pre(middleware.TrailMiddleware())

	e.Use(e_mid.RateLimiter(e_mid.NewRateLimiterMemoryStore(20)))
	// set up logger
	e.Use(e_mid.Logger())
	e.Use(e_mid.Recover())

	// Root route => handler
	e.HEAD("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, nil)
	})

	e.GET("/", func(c echo.Context) error {
		var resp = map[string]interface{}{
			"ApplicationName":     "Risigner Chat Server",
			"ApplicationOwner":    "Risigner Innovation Limited",
			"ApplicationVersion":  "1.0",
			"ApplicationEngineer": "Sam Ayo",
			"ApplicationStatus":   "running...",
		}
		return c.JSON(http.StatusOK, resp)
	})
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "healthy")
	})

	//set api endpoint
	api := e.Group("api/")

	api.Use(middleware.APIKeyMiddleware())
	SetupWebSocketRoutes(e, db, redisdb)
	//Run Server
	s := &http.Server{
		Addr:         ":" + string(config.AppConfig.PORT),
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}
	//s.SetKeepAlivesEnabled(false)
	e.HideBanner = true
	// Start server
	go func() {
		if err := e.StartServer(s); err != nil {
			log.Print(err.Error(), "shutting down the server")

		}

	}()
	log.Println("⚡️🚀 Risigner Chat Server::Started")
	log.Println("⚡️🚀 Risigner Chat Server::Running")
	ApiRoutes(api, db, redisdb)
	return e
}

// Stop - stop the echo server
func Stop(e *echo.Echo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	return nil
}
