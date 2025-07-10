package middleware

import (
	exp "chatsystem/internal/exceptions"
	"fmt"

	"github.com/labstack/echo/v4"
)

func SetHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Request().Header.Set("X-requested-With", "XMLHttpRequest")
		c.Request().Header.Set("Content-Type", "application/json")
		return next(c)
	}
}

func Recover(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				c.Error(err)
				exp.Loggers.System.Warn(fmt.Sprintf("Recovered from panic in endpoint: %v", r))
			}
		}()
		return next(c)
	}
}
