package errors

import "net/http"

var (
	ErrUserNotFound = NewAppError("USER_NOT_FOUND", "User not found", http.StatusNotFound,)
	ErrFailedToSaveUser =NewAppError("FAILED_TO_SAVE_USER", "Failed to save user", http.StatusInternalServerError)
	ErrFailedToCheckUser =NewAppError("FAILED_TO_CHECK_USER", "Failed to check user", http.StatusInternalServerError)
	ErrUserNotBelongToSpace =NewAppError("USER_NOT_BELONG_TO_SPACE", "User does not belong to space", http.StatusForbidden)
)
