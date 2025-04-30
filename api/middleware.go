package api

import (
	db "p2platform/db/sqlc"
	appErr "p2platform/errors"
	"p2platform/token"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func CookieAuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, err := ctx.Cookie("access_token")
		if err != nil {
			ctx.AbortWithStatusJSON(appErr.ErrCookieNotFound.Status, ErrorResponse(appErr.ErrCookieNotFound))
			return
		}
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(appErr.ErrInvalidCookie.Status, ErrorResponse(appErr.ErrInvalidCookie))
			log.Error().Str("Error: ", err.Error())
			return
		}
		//log.Info().Msg(fmt.Sprintf("Telegram id: %d\nTelegram username: %s", payload.TelegramId, payload.Username))
		// Установить payload.TelegramID в контекст
		ctx.Set("telegram_id", payload.TelegramId)
		ctx.Next()
		//ctx.Set("telegram_id", int64(86674601))
		//ctx.Next()
	}
}

func GetTelegramIDFromContext(ctx *gin.Context) (int64, bool) {
	telegramIdValue, exists := ctx.Get("telegram_id")
	if !exists {
		ctx.AbortWithStatusJSON(appErr.ErrUnauthorized.Status, ErrorResponse(appErr.ErrUnauthorized))
		return 0, false
	}
	telegramId, ok := telegramIdValue.(int64)
	if !ok {
		ctx.AbortWithStatusJSON(appErr.ErrUnauthorized.Status, ErrorResponse(appErr.ErrUnauthorized))
		return 0, false
	}
	return telegramId, true
}

func GetUserSpaces(ctx *gin.Context, store db.SQLStore) ([]uuid.UUID, error) {
	telegramIdValue, exists := ctx.Get("telegram_id")
	if !exists {
		ctx.JSON(appErr.ErrUnauthorized.Status, ErrorResponse(appErr.ErrUnauthorized))
		return nil, appErr.ErrUnauthorized
	}
	currentSpaceIds, err := store.GetSpaceIdByUserId(ctx, telegramIdValue.(int64))
	if err != nil {
		ctx.JSON(appErr.ErrInternalServer.Status, appErr.ErrInternalServer)
	}
	if len(currentSpaceIds) == 0 {

		ctx.JSON(appErr.ErrUserNotBelongToSpace.Status, appErr.ErrUserNotBelongToSpace)
		return nil, appErr.ErrUserNotBelongToSpace
	}
	return currentSpaceIds, nil
}
