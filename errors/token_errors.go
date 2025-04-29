package errors

import "net/http"


var (
	ErrTokenExpired = NewAppError("INVALID_TOKEN", "Token has been expired", http.StatusForbidden)
	ErrInvalidToken = NewAppError("INVALID_TOKEN", "Invalid token", http.StatusForbidden)
)