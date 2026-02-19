package handler

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	customerrors "github.com/sudo-hassan-zahid/go-api-server/internal/errors"
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
			errorHandler := customerrors.NewErrorHandler()
			statusCode, errorResponse := errorHandler.Handle(err)
			var ferr *fiber.Error
			if errors.As(err, &ferr) {
				statusCode = ferr.Code
				errorResponse.Error = ferr.Message
			}

			return ctx.Status(statusCode).JSON(errorResponse)
		},
	})
}
