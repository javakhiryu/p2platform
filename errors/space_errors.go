package errors

import "net/http"

var (
	ErrSpacesNotFound = NewAppError("SPACE_NOT_FOUND", "Space(s) not found", http.StatusNotFound)
	ErrSpaceNameExists = NewAppError("SPACE_NAME_EXISTS", "Space name already taken", http.StatusConflict)
	ErrIncorrectPassword = NewAppError("WRONG_PASSWORD", "Entered password is incorrect. Please, try again.", http.StatusForbidden)
	ErrUserAlreadyInSpace = NewAppError("USER_ALREADY_IN_SPACE", "User is already in space", http.StatusConflict)
)
