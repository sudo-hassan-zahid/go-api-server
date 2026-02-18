package database

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrUniqueViolation     = errors.New("unique constraint violation")
	ErrForeignKeyViolation = errors.New("foreign key constraint violation")
	ErrNotNullViolation    = errors.New("not null constraint violation")
	ErrCheckViolation      = errors.New("check constraint violation")
	ErrExclusionViolation  = errors.New("exclusion constraint violation")
)

type PostgresError struct {
	Message    string
	SQLState   string
	Constraint string
	Column     string
	Table      string
}

func ParsePostgresError(err error) *PostgresError {
	if err == nil {
		return nil
	}

	pgErr := &PostgresError{
		Message: err.Error(),
	}

	errStr := err.Error()

	if strings.Contains(errStr, "SQLSTATE") {
		parts := strings.Split(errStr, "SQLSTATE")
		if len(parts) > 1 {
			state := strings.TrimSpace(parts[1])
			if len(state) >= 5 {
				pgErr.SQLState = state[:5]
			}
		}
	}

	if strings.Contains(errStr, "constraint") {
		parts := strings.Split(errStr, "constraint")
		if len(parts) > 1 {
			constraint := strings.TrimSpace(parts[1])
			if strings.Contains(constraint, `"`) {
				start := strings.Index(constraint, `"`) + 1
				end := strings.Index(constraint[start:], `"`)
				if end > 0 {
					pgErr.Constraint = constraint[start : start+end]
				}
			}
		}
	}

	if strings.Contains(errStr, "column") {
		parts := strings.Split(errStr, "column")
		if len(parts) > 1 {
			column := strings.TrimSpace(parts[1])
			if strings.Contains(column, `"`) {
				start := strings.Index(column, `"`) + 1
				end := strings.Index(column[start:], `"`)
				if end > 0 {
					pgErr.Column = column[start : start+end]
				}
			}
		}
	}

	if strings.Contains(errStr, "table") {
		parts := strings.Split(errStr, "table")
		if len(parts) > 1 {
			table := strings.TrimSpace(parts[1])
			if strings.Contains(table, `"`) {
				start := strings.Index(table, `"`) + 1
				end := strings.Index(table[start:], `"`)
				if end > 0 {
					pgErr.Table = table[start : start+end]
				}
			}
		}
	}

	return pgErr
}

func IsUniqueViolation(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	pgErr := ParsePostgresError(err)
	if pgErr == nil {
		return false
	}

	return strings.Contains(pgErr.Message, "duplicate key value violates unique constraint") ||
		strings.Contains(pgErr.Message, "unique constraint") ||
		pgErr.SQLState == "23505"
}

func IsForeignKeyViolation(err error) bool {
	pgErr := ParsePostgresError(err)
	if pgErr == nil {
		return false
	}

	return strings.Contains(pgErr.Message, "foreign key constraint") ||
		strings.Contains(pgErr.Message, "violates foreign key constraint") ||
		pgErr.SQLState == "23503"
}

func IsNotNullViolation(err error) bool {
	pgErr := ParsePostgresError(err)
	if pgErr == nil {
		return false
	}

	return strings.Contains(pgErr.Message, "null value in column") ||
		strings.Contains(pgErr.Message, "not-null constraint") ||
		pgErr.SQLState == "23502"
}

func IsCheckViolation(err error) bool {
	pgErr := ParsePostgresError(err)
	if pgErr == nil {
		return false
	}

	return strings.Contains(pgErr.Message, "check constraint") ||
		strings.Contains(pgErr.Message, "violates check constraint") ||
		pgErr.SQLState == "23514"
}

func IsExclusionViolation(err error) bool {
	pgErr := ParsePostgresError(err)
	if pgErr == nil {
		return false
	}

	return strings.Contains(pgErr.Message, "exclusion constraint") ||
		strings.Contains(pgErr.Message, "violates exclusion constraint") ||
		pgErr.SQLState == "23P01"
}

func IsRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func GetConstraintName(err error) string {
	pgErr := ParsePostgresError(err)
	if pgErr == nil {
		return ""
	}
	return pgErr.Constraint
}

func GetColumnName(err error) string {
	pgErr := ParsePostgresError(err)
	if pgErr == nil {
		return ""
	}
	return pgErr.Column
}

func GetTableName(err error) string {
	pgErr := ParsePostgresError(err)
	if pgErr == nil {
		return ""
	}
	return pgErr.Table
}
