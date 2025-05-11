package db

import (
	"context"
	"database/sql"
	"errors"
	appErr "p2platform/errors"
	"p2platform/util"
)

type CreateUserTxParams struct {
	TelegramID int64  `json:"telegram_id"`
	TgUsername string `json:"tg_username"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	PhotoUrl   string `json:"photo_url"`
}

type CreateUserTxResult struct {
	User User `json:"user"`
}

func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		// 1. Проверяем существование пользователя
		_, err := q.GetUser(ctx, arg.TelegramID)
		if err == nil {
			// Пользователь уже существует
			return appErr.ErrUserAlreadyExists
		}

		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}

		// 2. Создаём нового пользователя
		user, err := q.CreateUser(ctx, CreateUserParams{
			TelegramID: arg.TelegramID,
			TgUsername: arg.TgUsername,
			FirstName:  util.ToPgText(arg.FirstName),
			LastName:   util.ToPgText(arg.LastName),
			PhotoUrl:   util.ToPgText(arg.PhotoUrl),
		})
		if err != nil {
			return err
		}

		result.User = user
		return nil
	})

	return result, err
}
