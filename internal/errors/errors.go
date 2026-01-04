package errors

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrBadRequest         = errors.New("bad request")
	ErrInternalServer     = errors.New("internal server error")
	ErrTokenInvalid       = errors.New("invalid or expired token")
)

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func SendError(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(fiber.Map{"error": message})
}

func HandleError(c *fiber.Ctx, err error) error {
	switch err {
	case ErrUserNotFound:
		return SendError(c, fiber.StatusNotFound, err.Error())
	case ErrInvalidCredentials:
		return SendError(c, fiber.StatusUnauthorized, err.Error())
	case ErrEmailAlreadyExists:
		return SendError(c, fiber.StatusBadRequest, err.Error())
	case ErrUnauthorized:
		return SendError(c, fiber.StatusUnauthorized, err.Error())
	case ErrForbidden:
		return SendError(c, fiber.StatusForbidden, err.Error())
	case ErrBadRequest:
		return SendError(c, fiber.StatusBadRequest, err.Error())
	default:
		return SendError(c, fiber.StatusInternalServerError, ErrInternalServer.Error())
	}
}
