package api

import (
	appErr "p2platform/errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CookieAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("telegram_id")
		if err != nil {
			ctx.AbortWithStatusJSON(appErr.ErrCookieNotFound.Status, ErrorResponse(appErr.ErrCookieNotFound))
			return
		}
		telegramID, err := strconv.ParseInt(cookie, 10, 64)
		if err != nil {
			ctx.AbortWithStatusJSON(appErr.ErrInvalidCookie.Status, ErrorResponse(appErr.ErrInvalidCookie))
			return
		}

		ctx.Set("telegram_id", telegramID)

		ctx.Next()
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
