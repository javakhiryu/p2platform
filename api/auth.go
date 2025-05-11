package api

import (
	"errors"
	"fmt"
	"net/http"
	"p2platform/auth"
	db "p2platform/db/sqlc"
	appErr "p2platform/errors"
	"p2platform/notification/kafka"
	"p2platform/notification/model"
	"p2platform/util"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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
		ctx.JSON(appErr.ErrInvalidPayload.Status, ErrorResponse(appErr.ErrInvalidPayload))
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
			ctx.JSON(appErr.ErrFailedToSaveUser.Status, appErr.ErrFailedToSaveUser)
			return
		}
	} else if err != nil {
		ctx.JSON(appErr.ErrFailedToCheckUser.Status, appErr.ErrFailedToCheckUser)
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(user.ID, user.Username)
	if err != nil {
		ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(err))
		log.Error().Str("Error:", err.Error())
		return
	}

	duration, err := time.ParseDuration(server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(err))
		log.Error().Str("Error:", err.Error())
		return
	}

	ctx.SetCookie(
		"access_token",
		accessToken,
		int(duration.Seconds()),
		"/",
		server.config.BaseURL,
		false,
		true,
	)

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

type getCurrentUserResponse struct {
	IsUserAuthorized bool `json:"is_user_authorized"`
}

func (server *Server) getCurrentUser(ctx *gin.Context) {
	var response getCurrentUserResponse
	response.IsUserAuthorized = true
	telegramID, ok := GetTelegramIDFromContext(ctx)
	if !ok {
		response.IsUserAuthorized = false
		return
	}
	_, err := server.store.GetUser(ctx, telegramID)
	if err != nil {
		if errors.Is(err, db.ErrNoRowsFound) {
			response.IsUserAuthorized = false
			ctx.JSON(http.StatusUnauthorized, response)
			return
		}
		log.Error().Err(err).Msg("error:")
		ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, response)
}

type AuthInitResponse struct {
	AuthUrl  string `json:"auth_url"`
	AuthCode string `json:"auth_code"`
	TTL      int    `json:"ttl"`
}

func (server *Server) initTelegramAuth(ctx *gin.Context) {
	authCode := uuid.New().String()
	expiresAt := time.Now().Add(time.Duration(server.config.TelegramAuthTTL.Minutes()))

	_, err := server.store.CreateTelegramAuthCode(ctx, db.CreateTelegramAuthCodeParams{
		AuthCode:  authCode,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		log.Error().Err(err).Msg("error:")
		ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(appErr.ErrInternalServer))
		return
	}
	authUrl := fmt.Sprintf("https://t.me/%s?start=%s", server.config.TelegramBotUsername, authCode)
	ctx.JSON(http.StatusOK, AuthInitResponse{
		AuthUrl:  authUrl,
		AuthCode: authCode,
		TTL:      int(server.config.TelegramAuthTTL.Seconds()),
	})
}

func (server *Server) handleTelegramWebhook(ctx *gin.Context) {
	var update tgbotapi.Update
	if err := ctx.BindJSON(&update); err != nil {
		ctx.JSON(appErr.ErrInvalidPayload.Status, ErrorResponse(appErr.ErrInvalidPayload))
		return
	}
	if update.Message.IsCommand() && update.Message.Command() == "start" {
		args := update.Message.CommandArguments()
		if strings.HasPrefix(args, "auth_") {
			authCode := strings.TrimPrefix(args, "auth_")
			err := server.store.ConfirmTelegramAuthCode(ctx, db.ConfirmTelegramAuthCodeParams{
				AuthCode:   authCode,
				TelegramID: util.ToPgInt(update.Message.From.ID),
				Status:     "confirmed",
			})
			if err != nil {
				log.Error().Err(err).Msg("error:")
				ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(appErr.ErrInternalServer))
				return
			}

			userTg := update.Message.From

			user, err := server.store.CreateUser(ctx, db.CreateUserParams{
				TelegramID: userTg.ID,
				TgUsername: userTg.UserName,
				PhotoUrl:   util.ToPgText(""),
				FirstName:  util.ToPgText(userTg.FirstName),
				LastName:   util.ToPgText(userTg.LastName),
			})
			if err != nil {
				log.Error().Err(err).Msg("error:")
				ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(appErr.ErrInternalServer))
				return
			}

			accessToken, err := server.tokenMaker.CreateToken(user.TelegramID, user.TgUsername)
			if err != nil {
				ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(err))
				log.Error().Str("Error:", err.Error())
				return
			}

			duration, err := time.ParseDuration(server.config.AccessTokenDuration)
			if err != nil {
				ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(err))
				log.Error().Str("Error:", err.Error())
				return
			}

			ctx.SetCookie(
				"access_token",
				accessToken,
				int(duration.Seconds()),
				"/",
				server.config.BaseURL,
				false,
				true,
			)

			_ = kafka.Publish(server.producer, "notifications", model.NotifictationMessage{
				TelegramId: update.Message.Chat.ID,
				Message:    fmt.Sprintf("✅ Авторизация успешна!"),
				EventType:  "auth_success",
			})
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

type AuthStatusRequest struct {
	Code string `form: "code", binding:"required,uuid"`
}

type AuthStatusResponse struct {
	Status     string `json:"status"`
	TelegramID int64  `json:"telegram_id, omitempty"`
}

func (server *Server) checkAuthStatus(ctx *gin.Context) {
	var req AuthStatusRequest
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(appErr.ErrInvalidQuery.Status, ErrorResponse(appErr.ErrInvalidQuery))
		return
	}
	code, err := server.store.GetTelegramAuthCode(ctx, req.Code)
	if err != nil {
		ctx.JSON(appErr.ErrInternalServer.Status, ErrorResponse(appErr.ErrInternalServer))
		return
	}
	if time.Now().After(code.ExpiresAt) {
		ctx.JSON(http.StatusOK, AuthStatusResponse{
			Status: "expired",
		})
		return
	}

	ctx.JSON(http.StatusOK, AuthStatusResponse{
		Status:     code.Status,
		TelegramID: code.TelegramID.Int64,
	})
}
