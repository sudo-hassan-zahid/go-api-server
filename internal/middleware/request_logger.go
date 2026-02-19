package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sudo-hassan-zahid/go-api-server/internal/logger"
)

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)
		logger.Log.Info().
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", c.Response().StatusCode()).
			Dur("latency", latency).
			Str("ip", c.IP()).
			Msg("http request")

		return err
	}
}
