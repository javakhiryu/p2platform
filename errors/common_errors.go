package errors

import "net/http"

var (
	ErrInternalServer = NewAppError("INTERNAL_SERVER_ERROR", "Something went wrong", http.StatusInternalServerError)
	ErrNotImplemented = NewAppError("NOT_IMPLEMENTED", "Not implemented", http.StatusNotImplemented)
	ErrInvalidPayload = NewAppError("INVALID_PAYLOAD", "Invalid payload", http.StatusBadRequest)
	ErrInvalidUri     = NewAppError("INVALID_URI", "Invalid URI", http.StatusBadRequest)
	ErrInvalidQuery   = NewAppError("INVALID_QUERY", "Invalid query", http.StatusBadRequest)
	ErrInvalidUUID    = NewAppError("INVALID_UUID", "Invalid UUID", http.StatusBadRequest)
)
