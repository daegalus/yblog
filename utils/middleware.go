package utils

import (
	"time"

	"github.com/caarlos0/log"
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
)

// Handler gin middleware handler
func Caarlos0Logger(message string) gin.HandlerFunc {
	var skip map[string]struct{}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			// Stop timer
			end := time.Now()
			latency := end.Sub(start)

			clientIP := c.ClientIP()
			method := c.Request.Method
			statusCode := c.Writer.Status()

			log.WithFields(log.Fields{
				"status_code": statusCode,
				"latency":     latency,
				"client_ip":   clientIP,
				"method":      method,
				"time":        time.Now().Format(time.RFC3339),
			}).Info(path)
		}
	}
}

func EchoCaarlos0Logger(message string) echo.MiddlewareFunc {
	var skip map[string]struct{}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Start timer
			start := time.Now()
			path := c.Request().URL.Path

			// Process request
			err := next(c)
			if err != nil {
				return err
			}

			// Log only when path is not being skipped
			if _, ok := skip[path]; !ok {
				// Stop timer
				end := time.Now()
				latency := end.Sub(start)

				clientIP := c.RealIP()
				method := c.Request().Method
				statusCode := c.Response().Status

				log.WithFields(log.Fields{
					"status_code": statusCode,
					"latency":     latency,
					"client_ip":   clientIP,
					"method":      method,
					"time":        time.Now().Format(time.RFC3339),
				}).Info(path)
			}
			return nil
		}
	}
}
