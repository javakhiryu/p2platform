package api

import (
	"fmt"
	appErr "p2platform/errors"
	"p2platform/token"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func CookieAuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//accessToken, err := ctx.Cookie("access_token")
		//if err != nil {
		//	ctx.AbortWithStatusJSON(appErr.ErrCookieNotFound.Status, ErrorResponse(appErr.ErrCookieNotFound))
		//	return
		//}
		payload, err := tokenMaker.VerifyToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZWxlZ3JhbV9pZCI6ODY2NzQ2MDEsInVzZXJuYW1lIjoiamF2YWtoaXJ5dSJ9.rOqzoEqkQ7kzmaiv6acKUtuNyJuduJmMw9z1hLUy8aw")
		if err != nil {
			ctx.AbortWithStatusJSON(appErr.ErrInvalidCookie.Status, ErrorResponse(appErr.ErrInvalidCookie))
			log.Error().Str("Error: ", err.Error())
			return
		}
		log.Info().Msg(fmt.Sprintf("Telegram id: %d\nTelegram username: %s", payload.TelegramId, payload.Username))
		// Установить payload.TelegramID в контекст
		ctx.Set("telegram_id", payload.TelegramId)
		ctx.Next()
		//ctx.Set("telegram_id", int64(86674601))
		//ctx.Next()
	}
}

func GetTelegramIDFromContext(ctx *gin.Context) (int64, bool) {
	tgIdValue, exists := ctx.Get("telegram_id")
	if !exists {
		ctx.JSON(appErr.ErrUnauthorized.Status, ErrorResponse(appErr.ErrUnauthorized))
		return 0, false
	}
	telegramId, ok := tgIdValue.(int64)
	if !ok {
		ctx.JSON(appErr.ErrUnauthorized.Status, ErrorResponse(appErr.ErrUnauthorized))
		return 0, false
	}
	return telegramId, true
}
