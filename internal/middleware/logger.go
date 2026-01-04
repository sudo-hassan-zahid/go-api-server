package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	appErrors "github.com/sudo-hassan-zahid/go-api-server/internal/errors"
	appLogger "github.com/sudo-hassan-zahid/go-api-server/internal/logger"
)

func ErrorLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		if err != nil {
			var status int
			var message string

			switch err {
			case appErrors.ErrBadRequest:
				status = fiber.StatusBadRequest
				message = err.Error()
			case appErrors.ErrUnauthorized:
				status = fiber.StatusUnauthorized
				message = err.Error()
			case appErrors.ErrForbidden:
				status = fiber.StatusForbidden
				message = err.Error()
			case appErrors.ErrInvalidCredentials:
				status = fiber.StatusUnauthorized
				message = err.Error()
			case appErrors.ErrEmailAlreadyExists:
				status = fiber.StatusBadRequest
				message = err.Error()
			case appErrors.ErrUserNotFound:
				status = fiber.StatusNotFound
				message = err.Error()
			default:
				status = fiber.StatusInternalServerError
				message = appErrors.ErrInternalServer.Error()
			}

			appLogger.Log.Error().
				Err(err).
				Str("method", c.Method()).
				Str("url", c.OriginalURL()).
				Int("status", status).
				Dur("latency", time.Since(start)).
				Msg("Request failed")

			return c.Status(status).JSON(fiber.Map{"error": message})
		}

		return nil
	}
}
