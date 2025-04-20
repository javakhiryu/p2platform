package errors

import "net/http"

var (
	ErrFailedToReleaseLockedAmount          = NewAppError("FAILED_TO_RELEASE_LOCKED_AMOUNT", "Failed to release locked amount", http.StatusInternalServerError)
	ErrLockedAmountsNotFound                = NewAppError("LOCKED_AMOUNTS_NOT_FOUND", "Locked amount(s) not found", http.StatusNotFound)
	ErrFailedToGetLockedAmountBySellRequest = NewAppError("FAILED_TO_GET_LOCKED_AMOUNT_BY_SELL_REQUEST", "Failed to get locked amount by sell request", http.StatusInternalServerError)
	ErrFailedToGetLockedAmountByBuyRequest  = NewAppError("FAILED_TO_GET_LOCKED_AMOUNT_BY_BUY_REQUEST", "Failed to get locked amount by buy request", http.StatusInternalServerError)
	ErrFailedToGetLockedAmountByUser        = NewAppError("FAILED_TO_GET_LOCKED_AMOUNT_BY_USER", "Failed to get locked amount by user", http.StatusInternalServerError)
	ErrFailedToCreateLockedAmount           = NewAppError("FAILED_TO_CREATE_LOCKED_AMOUNT", "Failed to create locked amount", http.StatusInternalServerError)
)
