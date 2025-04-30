package errors

import "net/http"

var (
	ErrSpacesNotFound = NewAppError("SPACE_NOT_FOUND", "Space(s) not found", http.StatusNotFound)
	ErrSpaceNameExists = NewAppError("SPACE_NAME_EXISTS", "Space name already taken", http.StatusConflict)
)
