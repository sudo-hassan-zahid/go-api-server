package errors

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type AppError struct {
	Err        error
	HTTPStatus int
	Message    string
}

func (e *AppError) Error() string {
	return e.Err.Error()
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(err error, status int, message string) *AppError {
	return &AppError{
		Err:        err,
		HTTPStatus: status,
		Message:    message,
	}
}

// App errors
var (
	ErrTokenInvalid       = New(errors.New("invalid token"), fiber.StatusUnauthorized, "invalid token")
	ErrUserNotFound       = New(errors.New("user not found"), fiber.StatusNotFound, "user not found")
	ErrUserAlreadyExists  = New(errors.New("user already exists"), fiber.StatusConflict, "user already exists")
	ErrInvalidCredentials = New(errors.New("invalid credentials"), fiber.StatusUnauthorized, "invalid credentials")
	ErrUnauthorized       = New(errors.New("unauthorized"), fiber.StatusUnauthorized, "unauthorized")
	ErrValidationFailed   = New(errors.New("validation failed"), fiber.StatusBadRequest, "validation failed")
	ErrServiceUnavailable = New(errors.New("service temporarily unavailable"), fiber.StatusServiceUnavailable, "service temporarily unavailable")
)

// DB errors
var (
	ErrDuplicateKey        = New(errors.New("duplicate key"), fiber.StatusConflict, "resource already exists")
	ErrForeignKeyViolation = New(errors.New("foreign key violation"), fiber.StatusBadRequest, "referenced resource does not exist")
	ErrNotNullViolation    = New(errors.New("not null violation"), fiber.StatusBadRequest, "required field cannot be empty")
	ErrRecordNotFound      = New(errors.New("record not found"), fiber.StatusNotFound, "resource not found")
	ErrDatabaseTimeout     = New(errors.New("database timeout"), fiber.StatusRequestTimeout, "database operation timed out")
	ErrConnectionFailed    = New(errors.New("connection failed"), fiber.StatusServiceUnavailable, "database temporarily unavailable")
	ErrDeadlockDetected    = New(errors.New("deadlock detected"), fiber.StatusConflict, "database conflict, please try again")
)
