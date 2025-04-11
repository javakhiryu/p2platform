package db

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ForeighnKeyViolation = "23503"
	UniqueViolation      = "23505"
)

var ErrNoRowsFound = pgx.ErrNoRows

var ErrUniqueViolation = &pgconn.PgError{
	Code: UniqueViolation,
}

var ErrForeighnKeyViolation = &pgconn.PgError{
	Code: ForeighnKeyViolation,
}

var ErrSellRequestAlreadyDeleted = errors.New("sell request has been already deleted")
var ErrSEllRequestNotFound = errors.New("sell request not found")

func ErrCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}
