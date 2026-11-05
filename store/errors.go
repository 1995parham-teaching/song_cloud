package store

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	// ErrNotFound is returned when a referenced record (playlist, user or song) does not exist.
	ErrNotFound = errors.New("not found")
	// ErrForbidden is returned when the user has no access to the requested record.
	ErrForbidden = errors.New("forbidden")
	// ErrDuplicate is returned when a record with the same identity already exists.
	ErrDuplicate = errors.New("already exists")
	// ErrInvalid is returned when a request violates a business rule.
	ErrInvalid = errors.New("invalid request")
)

// PostgreSQL error codes, see https://www.postgresql.org/docs/current/errcodes-appendix.html.
const (
	codeForeignKeyViolation = "23503"
	codeUniqueViolation     = "23505"
)

// hasCode reports whether err is a PostgreSQL error with the given code.
func hasCode(err error, code string) bool {
	pgErr, ok := errors.AsType[*pgconn.PgError](err)

	return ok && pgErr.Code == code
}
