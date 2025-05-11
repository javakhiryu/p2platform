package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	CreateBuyRequestTx(ctx context.Context, arg CreateBuyRequestTxParams) (CreateBuyRequestTxResult, error)
	DeleteSellRequestTx(ctx context.Context, sellRequestId int32) (bool, error)
	CloseBuyRequestTx(ctx context.Context, arg CloseBuyRequestTxParams) (CloseBuyRequestTxResult, error)
	DeleteBuyRequestTx(ctx context.Context, buyReqID uuid.UUID) (bool, error)
	GetSellRequestTx(ctx context.Context, sellReqID int32, requesterTelegramID int64) (GetSellRequestTxResult, error)
	ListSellRequeststTx(ctx context.Context, params ListSellRequeststTxParams, requesterTelegramID int64) (ListSellRequeststTxResults, error)
	ReleaseLockedAmountTx(ctx context.Context, buyReqID uuid.UUID) (result ReleaseLockedAmountTxResult, err error)
	ListMySellRequeststTx(ctx context.Context, params ListMySellRequestsTxParams) (ListSellRequeststTxResults, error)
	CreateSpaceTx(ctx context.Context, arg CreateSpaceTxParams) (CreateSpaceTxResult, error)
	GetSpaceTx(ctx context.Context, spaceId uuid.UUID, requesterTelegramID int64) (GetSpaceTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
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
