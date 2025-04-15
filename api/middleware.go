package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CookieAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("telegram_id")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: cookie not found"})
			return
		}
		telegramID, err := strconv.ParseInt(cookie, 10, 64)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: invalid cookie"})
			return
		}

		ctx.Set("telegram_id", telegramID)

		ctx.Next()
	}
}

func GetTelegramIDFromContext(ctx *gin.Context) (int64, bool) {
	tgIdValue, exists := ctx.Get("telegram_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "telegram_id not found in context"})
		return 0, false
	}
	telegramId, ok := tgIdValue.(int64)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "telegram_id has invalid type"})
		return 0, false
	}
	return telegramId, true
}
