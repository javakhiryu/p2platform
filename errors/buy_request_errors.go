package errors

import "net/http"

var (
	ErrFailedToCreateBuyRequest = NewAppError("CREATE_BUY_REQUEST_FAILED", "Failed to create buy request", http.StatusInternalServerError)
	ErrFailedToCloseBuyRequests = NewAppError("CLOSE_BUY_REQUESTS_FAILED", "Failed to close buy request(s)", http.StatusInternalServerError)
	ErrBuyRequestsNotFound       = NewAppError("BUY_REQUEST_NOT_FOUND", "Buy request(s) not found", http.StatusNotFound)
	ErrBuyRequestDeleted        = NewAppError("BUY_REQUEST_DELETED", "Buy request has been deleted", http.StatusGone)
	ErrBuyRequestAlreadyDeleted = NewAppError("BUY_REQUEST_ALREADY_DELETED", "Buy request has already been deleted", http.StatusConflict)
	ErrFailedToGetBuyRequests   = NewAppError("FAILED_TO_GET_BUY_REQUESTS", "Failed to get buy requests", http.StatusInternalServerError)
	ErrFailedToDeleteBuyRequest = NewAppError("DELETE_BUY_REQUEST_FAILED", "Failed to delete buy request", http.StatusInternalServerError)
	ErrInsuficientTotalFunds    = NewAppError("INSUFFICIENT_FUNDS", "Not enough funds to create buy request", http.StatusBadRequest)
	ErrInsuficientCardFunds     = NewAppError("INSUFFICIENT_CARD_FUNDS", "Not enough card funds to create buy request", http.StatusBadRequest)
	ErrInsuficientCashFunds     = NewAppError("INSUFFICIENT_CASH_FUNDS", "Not enough cash funds to create buy request", http.StatusBadRequest)
	ErrTotalBuyAmountMismatch   = NewAppError("BUY_AMOUNT_MISMATCH", "Sum of cash and card amounts must equal total amount", http.StatusBadRequest)
	ErrNotBuyRequestOwner = NewAppError("NOT_BUY_REQUEST_OWNER", "You are not the owner of this buy request", http.StatusForbidden)
)
