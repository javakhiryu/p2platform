package db

import (
	"context"
	"errors"
	appErr "p2platform/errors"
	"p2platform/util"

	"github.com/google/uuid"
)

type CreateSpaceTxParams struct {
	SpaceID        uuid.UUID `json:"space_id"`
	SpaceName      string    `json:"space_name"`
	HashedPassword string    `json:"hashed_password"`
	Description    string    `json:"description"`
	CreatorID      int64     `json:"creator_id"`
}

type CreateSpaceTxResult struct {
	SpaceID     uuid.UUID `json:"space_id"`
	SpaceName   string    `json:"space_name"`
	User        User      `json:"user"`
	Description string    `json:"description"`
}

func (store *SQLStore) CreateSpaceTx(ctx context.Context, arg CreateSpaceTxParams) (CreateSpaceTxResult, error) {
	var result CreateSpaceTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		user, err := q.GetUser(ctx, arg.CreatorID)
		if err != nil {
			if errors.Is(err, ErrNoRowsFound) {
				return appErr.ErrUserNotFound
			}
			return appErr.ErrInternalServer
		}
		argCreateSpace := CreateSpaceParams{
			SpaceID:        arg.SpaceID,
			SpaceName:      arg.SpaceName,
			HashedPassword: arg.HashedPassword,
			Description:    arg.Description,
			CreatorID:      util.ToPgInt(user.TelegramID),
		}
		space, err := q.CreateSpace(ctx, argCreateSpace)
		if err != nil {
			if ErrCode(err) == UniqueViolation {
				return appErr.ErrSpaceNameExists
			}
			return appErr.ErrInternalServer
		}
		//admin automatically becomes the space member
		_, err = q.AddSpaceMember(ctx, AddSpaceMemberParams{
			SpaceID:  space.SpaceID,
			UserID:   user.TelegramID,
			Username: user.TgUsername,
		})
		if err != nil {
			return appErr.ErrInternalServer
		}
		result = CreateSpaceTxResult{
			SpaceID:     space.SpaceID,
			SpaceName:   space.SpaceName,
			User:        user,
			Description: space.Description,
		}
		return nil
	})
	return result, err
}
