package logger

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

// EchoLogger returns a middleware that logs HTTP requests using zerolog
// func EchoLogger() echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			start := time.Now()

// 			// Process request
// 			err := next(c)

// 			// Log request details
// 			duration := time.Since(start)
// 			event := log.Info().
// 				Str("remote_ip", c.Request().RemoteAddr).
// 				Str("host", c.Request().Host).
// 				Str("method", c.Request().Method).
// 				Str("uri", c.Request().URL.Path).
// 				Str("user_agent", c.Request().UserAgent()).
// 				Int("status", c.Response().Status).
// 				Dur("duration", duration)

// 			if err != nil {
// 				event.Err(err)
// 			}

// 			event.Msg("http request")

// 			return err
// 		}
// 	}
// }

func EchoLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		// LogLatency instructs logger to record duration it took to execute rest of the handler chain (next(c) call).
		LogLatency: true,
		// LogProtocol instructs logger to extract request protocol (i.e. `HTTP/1.1` or `HTTP/2`)
		LogProtocol: true,
		// LogRemoteIP instructs logger to extract request remote IP. See `echo.Context.RealIP()` for implementation details.
		LogRemoteIP: true,
		// LogHost instructs logger to extract request host value (i.e. `example.com`)
		LogHost: true,
		// LogMethod instructs logger to extract request method value (i.e. `GET` etc)
		LogMethod: true,
		// LogURI instructs logger to extract request URI (i.e. `/list?lang=en&page=1`)
		LogURI: true,
		// LogURIPath instructs logger to extract request URI path part (i.e. `/list`)
		LogURIPath: true,
		// LogRoutePath instructs logger to extract route path part to which request was matched to (i.e. `/user/:id`)
		LogRoutePath: true,
		// LogRequestID instructs logger to extract request ID from request `X-Request-ID` header or response if request did not have value.
		LogRequestID: true,
		// LogReferer instructs logger to extract request referer values.
		LogReferer: true,
		// LogUserAgent instructs logger to extract request user agent values.
		LogUserAgent: true,
		// LogStatus instructs logger to extract response status code. If handler chain returns an echo.HTTPError,
		// the status code is extracted from the echo.HTTPError returned
		LogStatus: true,
		// LogError instructs logger to extract error returned from executed handler chain.
		LogError: true,
		// LogContentLength instructs logger to extract content length header value. Note: this value could be different from
		// actual request body size as it could be spoofed etc.
		LogContentLength: true,
		// LogResponseSize instructs logger to extract response content length value. Note: when used with Gzip middleware
		// this value may not be always correct.
		LogResponseSize: true,
		LogHeaders:      []string{},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Info().
				Str("remote_ip", v.RemoteIP).
				Str("host", v.Host).
				Str("method", v.Method).
				Str("uri", v.URI).
				Str("user_agent", v.UserAgent).
				Int("status", v.Status).
				Dur("duration", v.Latency).
				Int64("bytes_in", c.Request().ContentLength).
				Int64("bytes_out", v.ResponseSize).
				Str("request_id", v.RequestID).
				Msg("http request")

			return nil
		},
	})
}
