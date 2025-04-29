package errors

import "net/http"

var (
	ErrUnauthorized = NewAppError("UNAUTHORIZED", "Unauthorized access",http.StatusUnauthorized)
	ErrForbidden =NewAppError("ACCESS_DENIED", "Access denied", http.StatusForbidden)
	ErrCookieNotFound =NewAppError("INVALID_COOKIE", "Cookie not found", http.StatusUnauthorized)
	ErrInvalidCookie =NewAppError("INVALID_COOKIE", "Invalid cookie", http.StatusUnauthorized)
	ErrFailedToHashPassword =NewAppError("FAILED_TO_HASH_PASSWORD", "Failed to hash password", http.StatusInternalServerError)
)