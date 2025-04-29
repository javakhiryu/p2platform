package util

import (
	appErr "p2platform/errors"

	"golang.org/x/crypto/bcrypt"
)

func HashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", appErr.ErrFailedToHashPassword
	}
	return string(hashedPassword), nil
}

func CheckPassword(password string, hashedPassword string) error{
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
