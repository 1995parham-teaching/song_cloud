package db

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
)

// PgError represents a detailed PostgreSQL error response.
type PgError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
	Hint    string `json:"hint,omitempty"`
	Where   string `json:"where,omitempty"`
}

// HandleError converts a database error to an appropriate HTTP error with details.
func HandleError(err error) *echo.HTTPError {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		status := http.StatusInternalServerError

		// Map PostgreSQL error classes to HTTP status codes
		switch pgErr.Code[:2] {
		case "23": // Integrity constraint violation
			status = http.StatusConflict
		case "P0": // PL/pgSQL error (user-defined exceptions)
			status = http.StatusBadRequest
		case "22": // Data exception
			status = http.StatusBadRequest
		case "42": // Syntax error or access rule violation
			status = http.StatusBadRequest
		}

		return echo.NewHTTPError(status, PgError{
			Code:    pgErr.Code,
			Message: pgErr.Message,
			Detail:  pgErr.Detail,
			Hint:    pgErr.Hint,
			Where:   pgErr.Where,
		})
	}

	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}
