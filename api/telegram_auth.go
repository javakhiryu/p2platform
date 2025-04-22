package api

import (
	"fmt"
	"net/http"
	"p2platform/auth"
	db "p2platform/db/sqlc"
	"p2platform/errors"
	appErr "p2platform/errors"
	"p2platform/util"

	"github.com/gin-gonic/gin"
)

type user struct {
	ID        int64  `json:"id" binding:"required"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	PhotoUrl  string `json:"photo_url"`
	Hash      string `json:"hash" binding:"required"`
	AuthDate  int64  `json:"auth_date" binding:"required"`
}

func (server *Server) telegramAuth(ctx *gin.Context) {
	user := user{}
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(errors.ErrInvalidPayload.Status, ErrorResponse(errors.ErrInvalidPayload))
		return
	}

	ok := auth.VerifyTelegramAuth(map[string]string{
		"id":         fmt.Sprint(user.ID),
		"username":   user.Username,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"photo_url":  user.PhotoUrl,
		"auth_date":  fmt.Sprint(user.AuthDate),
	}, user.Hash, server.config.TelegramBotToken)
	if !ok {
		ctx.JSON(appErr.ErrUnauthorized.Status, ErrorResponse(appErr.ErrUnauthorized))
		return
	}

	_, err := server.store.GetUser(ctx, user.ID)
	if err == db.ErrNoRowsFound {
		_, err := server.store.CreateUser(ctx, db.CreateUserParams{
			TelegramID: user.ID,
			TgUsername: user.Username,
			PhotoUrl:   util.ToPgText(user.PhotoUrl),
			FirstName:  util.ToPgText(user.FirstName),
			LastName:   util.ToPgText(user.LastName),
		})
		if err != nil {
			ctx.JSON(errors.ErrFailedToSaveUser.Status, errors.ErrFailedToSaveUser)
			return
		}
	} else if err != nil {
		ctx.JSON(errors.ErrFailedToCheckUser.Status, errors.ErrFailedToCheckUser)
		return
	}

	ctx.SetCookie(
		"telegram_id",
		fmt.Sprint(user.ID),
		3600*24*30,
		"/",
		"",
		true,
		true,
	)

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}
