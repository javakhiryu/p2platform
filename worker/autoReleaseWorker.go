package worker

import (
	"context"
	"fmt"
	db "p2platform/db/sqlc"
	"time"

	"github.com/rs/zerolog/log"
)

type AutoReleaseWorker struct {
	store  db.Store
	ticker *time.Ticker
	ctx    context.Context
}

func NewAutoReleaseWorker(store db.Store, interval time.Duration) *AutoReleaseWorker {
	return &AutoReleaseWorker{
		store:  store,
		ticker: time.NewTicker(interval),
		ctx:    context.Background(),
	}
}

func (w *AutoReleaseWorker) Start() {
	logger := log.Info()
	logger.Msg("AutoReleaseWorker started...")
	go func() {
		for {
			select {
			case <-w.ticker.C:
				w.processExpiredBuyRequests()
			case <-w.ctx.Done():
				log.Info().Msg("AutoReleaseWorker stopped")
				return
			}
		}
	}()
}

func (w *AutoReleaseWorker) Stop() {
	w.ticker.Stop()
}

func (w *AutoReleaseWorker) processExpiredBuyRequests() {
	expired, err := w.store.ListExpiredBuyRequests(w.ctx)
	if err != nil {
		log.Err(err).Msg("Failed to fetch expired buy requests")
		return
	}
	for _, req := range expired {
		result, err := w.store.ReleaseLockedAmountTx(w.ctx, req.BuyReqID)
		if err != nil {
			log.Err(err).Msg(fmt.Sprintf("Failed to release locked amount for buy_request: %v", req.BuyReqID))
		} else {
			log.Info().Msg(fmt.Sprintf("Auto released locked amount for buy_request %v", req.BuyReqID))
		}
		log.Info().
			Interface("release_result", result).
			Msg("ReleaseLockedAmountTx completed")

	}
}
