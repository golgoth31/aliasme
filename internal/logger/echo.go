package logger

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// EchoLogger returns a middleware that logs HTTP requests using zerolog
func EchoLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Process request
			err := next(c)

			// Log request details
			duration := time.Since(start)
			event := log.Info().
				Str("remote_ip", c.Request().RemoteAddr).
				Str("host", c.Request().Host).
				Str("method", c.Request().Method).
				Str("uri", c.Request().URL.Path).
				Str("user_agent", c.Request().UserAgent()).
				Int("status", c.Response().Status).
				Dur("duration", duration)

			if err != nil {
				event.Err(err)
			}

			event.Msg("http request")

			return err
		}
	}
}
