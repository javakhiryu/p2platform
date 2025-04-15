package api

import (
	"database/sql"
	"fmt"
	"net/http"
	db "p2platform/db/sqlc"
	"p2platform/auth"
	"p2platform/util"

	"github.com/gin-gonic/gin"
)

type user struct {
	ID        int64  `json:"id" binding:"required"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Hash      string `json:"hash" binding:"required"`
	AuthDate  int64  `json:"auth_date" binding:"required"`
}

func (s *Server) telegramAuth(ctx *gin.Context) {
	user := user{}
	config, err := util.LoadConfig(".")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load config"})
	}
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	ok := auth.VerifyTelegramAuth(map[string]string{
		"id":         fmt.Sprint(user.ID),
		"username":   user.Username,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"auth_date":  fmt.Sprint(user.AuthDate),
	}, user.Hash, config.TelegramBotToken)

	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid hash"})
		return
	}

	_, err = s.store.GetUser(ctx, user.ID)
	if err == sql.ErrNoRows {
		_, err := s.store.CreateUser(ctx, db.CreateUserParams{
			TelegramID: user.ID,
			TgUsername: user.Username,
			FirstName:  util.ToPgText(user.FirstName),
			LastName:   util.ToPgText(user.LastName),
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save user"})
			return
		}
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check user"})
		return
	}

	ctx.SetCookie(
		"telegram_id",
		fmt.Sprint(user.ID),
		3600*24*30,
		"/",
		"",
		false,
		true,
	)

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}
