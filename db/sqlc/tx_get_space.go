package db

import (
	appErr "p2platform/errors"

	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type GetSpaceTxResult struct {
	SpaceID      uuid.UUID `json:"space_id"`
	SpaceName    string    `json:"space_name"`
	Description  string    `json:"description"`
	Creator      User      `json:"creator"`
	MembersCount int64     `json:"members_count"`
}

func (store *SQLStore) GetSpaceTx(ctx context.Context, spaceId uuid.UUID, requesterTelegramID int64) (GetSpaceTxResult, error) {
	var result GetSpaceTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		arg := IsUserInSpaceParams{
			UserID:  requesterTelegramID,
			SpaceID: spaceId,
		}
		isUserInSpace, err := q.IsUserInSpace(ctx, arg)
		if err != nil {
			return appErr.ErrForbidden
		}
		if !isUserInSpace {
			return appErr.ErrForbidden
		}

		space, err := q.GetSpaceBySpaceId(ctx, spaceId)
		if err != nil {
			if err == ErrNoRowsFound{
				return appErr.ErrSpacesNotFound
			}
			return appErr.ErrInternalServer
		}
		creator, err := q.GetUser(ctx, space.CreatorID.Int64)
		if err != nil {
			if err == ErrNoRowsFound{
				return appErr.ErrUserNotFound
			}
			return appErr.ErrInternalServer
		}
		countOfSpaceMembers, err :=q.GetSpaceMembersCount(ctx, spaceId)
		if err != nil {
			return appErr.ErrInternalServer
		}
		result.SpaceID = space.SpaceID
		result.SpaceName = space.SpaceName
		result.Description = space.Description
		result.Creator = creator
		result.MembersCount = countOfSpaceMembers

		return nil
	})
	return result, err
}
