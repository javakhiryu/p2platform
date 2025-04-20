package errors

import "net/http"

var (
	ErrSellAmountMismatch = NewAppError(
		"SELL_AMOUNT_MISMATCH",
		"Sum of cash and card amounts must equal total amount",
		http.StatusBadRequest,
	)

	ErrSellRequestNotFound = NewAppError(
		"SELL_REQUEST_NOT_FOUND",
		"Sell request not found",
		http.StatusNotFound,
	)

	ErrSellRequestAlreadyDeleted = NewAppError(
		"SELL_REQUEST_ALREADY_DELETED",
		"Sell request has already been deleted",
		http.StatusConflict,
	)

	ErrSellRequestDeleted = NewAppError(
		"SELL_REQUEST_DELETED",
		"Sell request has been deleted",
		http.StatusGone,
	)
	ErrSellRequestIsNotActual = NewAppError(
		"SELL_REQUEST_IS_NOT_ACTUAL",
		"Sell request is not actual",
		http.StatusGone,
	)

	ErrNotSellRequestOwner = NewAppError(
		"NOT_SELL_REQUEST_OWNER",
		"You are not the owner of this sell request",
		http.StatusForbidden,
	)

	ErrSellRequestHasBuyRequests = NewAppError(
		"SELL_REQUEST_HAS_BUY_REQUESTS",
		"Sell request can't be updated because it has active buy requests",
		http.StatusForbidden,
	)

	ErrFailedToGetSellRequests = NewAppError(
		"GET_SELL_REQUESTS_FAILED",
		"Failed to get sell request(s)",
		http.StatusInternalServerError,
	)

	ErrFailedToCreateSellRequest = NewAppError(
		"CREATE_SELL_REQUEST_FAILED",
		"Failed to create sell request",
		http.StatusInternalServerError,
	)

	ErrFailedToUpdateSellRequest = NewAppError(
		"UPDATE_SELL_REQUEST_FAILED",
		"Failed to update sell request",
		http.StatusInternalServerError,
	)

	ErrFailedToDeleteSellRequest = NewAppError(
		"DELETE_SELL_REQUEST_FAILED",
		"Failed to delete sell request",
		http.StatusInternalServerError,
	)
	ErrFailedToCloseSellRequest = NewAppError(
		"FAILED_TO_CLOSE_SELL_REQUEST",
		"Failed to close related sell request(s)",
		http.StatusInternalServerError,
	)
	ErrFailedToOpenSellRequest = NewAppError(
		"FAILED_TO_OPEN_SELL_REQUEST",
		"Failed to open related sell request(s)",
		http.StatusInternalServerError,
	)
)
