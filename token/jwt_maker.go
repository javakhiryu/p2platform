package token

import (
	"fmt"
	appErr "p2platform/errors"

	"github.com/golang-jwt/jwt/v5"
)

const (
	minSecretKeySize = 32
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(telegramId int64, username string) (string, error) {
	claims := jwt.MapClaims{
		"telegram_id": telegramId,
		"username":    username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(tokenStr string) (*Payload, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, appErr.ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	})
	if err != nil {
		return nil, appErr.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, appErr.ErrInvalidToken
	}

	telegramID, ok := claims["telegram_id"].(float64)
	username, ok2 := claims["username"].(string)
	if !ok || !ok2 {
		return nil, appErr.ErrInvalidToken
	}

	return &Payload{
		TelegramId: int64(telegramID),
		Username:   username,
	}, nil
}
