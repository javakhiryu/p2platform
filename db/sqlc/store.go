package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	CreateBuyRequestTx(ctx context.Context, arg CreateBuyRequestTxParams) (CreateBuyRequestTxResult, error)
	DeleteSellRequestTx(ctx context.Context, sellRequestId int32) (bool, error)
	CloseBuyRequestTx(ctx context.Context, arg CloseBuyRequestTxParams) (CloseBuyRequestTxResult, error)
}

// Store
type SQLStore struct {
	*Queries
	connPool *pgxpool.Pool
}

func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
