package middleware

import (
	"chatsystem/internal/config"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// This enables us interact with the Frontend Frameworks
func CORSMiddleware() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowCredentials: true,
	})
}

func TrailMiddleware() echo.MiddlewareFunc {
	return middleware.RemoveTrailingSlash()
}

func APIKeyMiddleware() echo.MiddlewareFunc {
	expectedKey := config.AppConfig.APIKey

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			apiKey := c.Request().Header.Get("x-api-key")
			if apiKey == "" || apiKey != expectedKey {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing or invalid API key")
			}
			return next(c)
		}
	}
}
