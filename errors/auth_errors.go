package errors

import "net/http"

var (
	ErrUnauthorized = NewAppError("UNAUTHORIZED", "Unauthorized access",http.StatusUnauthorized)
	ErrForbidden =NewAppError("ACCESS_DENIED", "Access denied", http.StatusForbidden)
	ErrCookieNotFound =NewAppError("INVALID_COOKIE", "Cookie not found", http.StatusUnauthorized)
	ErrInvalidCookie =NewAppError("INVALID_COOKIE", "Invalid cookie", http.StatusUnauthorized)
)