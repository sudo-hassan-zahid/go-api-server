package handler

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sudo-hassan-zahid/go-api-server/internal/domain"
	"github.com/sudo-hassan-zahid/go-api-server/internal/logger"
)

func NewApp() *fiber.App {
	return fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		BodyLimit:    50 * 1024 * 1024,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			logger.Log.Error().Err(err).Msg("Failed to handle request")

			message := err.Error()

			var code int
			var ferr *fiber.Error

			if errors.As(err, &ferr) {
				code = ferr.Code
			} else {
				switch {
				case errors.Is(err, domain.ErrUserNotFound):
					code = fiber.StatusNotFound
				case errors.Is(err, domain.ErrUserAlreadyExists):
					code = fiber.StatusConflict
				case errors.Is(err, domain.ErrForeignKeyViolation):
					code = fiber.StatusBadRequest
				case errors.Is(err, domain.ErrNotNullViolation):
					code = fiber.StatusBadRequest
				case errors.Is(err, domain.ErrCheckViolation):
					code = fiber.StatusBadRequest
				case errors.Is(err, domain.ErrInvalidCredentials):
					code = fiber.StatusUnauthorized
				case errors.Is(err, domain.ErrUnauthorized):
					code = fiber.StatusUnauthorized
				case errors.Is(err, domain.ErrTokenInvalid):
					code = fiber.StatusUnauthorized
				default:
					code = fiber.StatusInternalServerError
					message = "internal server error"
				}
			}

			return ctx.Status(code).JSON(fiber.Map{
				"success": false,
				"error":   message,
			})
		},
	})
}
