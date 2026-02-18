package errors

import (
	"errors"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

func ParseDatabaseError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrRecordNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return ErrDuplicateKey
		case "23503":
			return ErrForeignKeyViolation
		case "23502":
		case "40P01":
			return ErrDeadlockDetected
		}
	}

	return err
}
