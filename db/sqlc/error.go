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
var ErrSellRequestNotFound = errors.New("sell request(s) not found")
var BuyRequestNotFoundOrDeleted = errors.New("buy request not found or was already deleted")

func ErrCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}
