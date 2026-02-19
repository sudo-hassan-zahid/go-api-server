package errors

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type ErrorHandler struct{}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

func (h *ErrorHandler) Handle(err error) (int, ErrorResponse) {
	if err == nil {
		return fiber.StatusOK, ErrorResponse{
			Success: true,
		}
	}

	err = ParseDatabaseError(err)

	log.Println("ERROR:", err)

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.HTTPStatus, ErrorResponse{
			Success: false,
			Error:   appErr.Message,
		}
	}

	return fiber.StatusInternalServerError, ErrorResponse{
		Success: false,
		Error:   "internal server error",
	}
}
