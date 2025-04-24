package db

import (
	"context"
	"database/sql"
	"p2platform/util"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func CreateRandomBuyRequest(t *testing.T, sellRequest SellRequest, user User) BuyRequest {
	arg := CreateBuyRequestParams{
		BuyReqID:        uuid.New(),
		SellReqID:       sellRequest.SellReqID,
		BuyTotalAmount:  sellRequest.SellTotalAmount,
		TelegramID:      user.TelegramID,
		TgUsername:      util.RandomTgUsername(),
		BuyByCard:       sellRequest.SellByCard,
		BuyAmountByCard: sellRequest.SellAmountByCard,
		BuyByCash:       sellRequest.SellByCash,
		BuyAmountByCash: sellRequest.SellAmountByCash,
	}
	buyRequest, err := testStore.CreateBuyRequest(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, buyRequest)
	require.Equal(t, sellRequest.SellReqID, buyRequest.SellReqID)
	return buyRequest
}

func TestCreateBuyRequest(t *testing.T) {
	CreateRandomBuyRequest(t, CreateRandomSellRequest(t, CreateRandomUser(t)), CreateRandomUser(t))
}

func TestGetBuyRequest(t *testing.T) {
	buyRequest := CreateRandomBuyRequest(t, CreateRandomSellRequest(t, CreateRandomUser(t)), CreateRandomUser(t))
	buyRequest2, err := testStore.GetBuyRequestById(context.Background(), buyRequest.BuyReqID)
	require.NoError(t, err)
	require.NotEmpty(t, buyRequest2)
	require.Equal(t, buyRequest.BuyReqID, buyRequest2.BuyReqID)
}

func TestListBuyRequests(t *testing.T) {
	sellRequest := CreateRandomSellRequest(t, CreateRandomUser(t))
	user := CreateRandomUser(t)
	var buyRequests []BuyRequest
	for i := 0; i < 5; i++ {
		buyRequests = append(buyRequests, CreateRandomBuyRequest(t, sellRequest, user))
	}
	arg := ListBuyRequestsParams{
		SellReqID: sellRequest.SellReqID,
		Limit:     5,
		Offset:    0,
	}
	buyRequests, err := testStore.ListBuyRequests(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, buyRequests)
	require.Len(t, buyRequests, 5)
	for _, br := range buyRequests {
		require.Equal(t, buyRequests[0].SellReqID, br.SellReqID)
	}
}

func TestListBuyReuqestsByTelegramId(t *testing.T) {
	sellRequest := CreateRandomSellRequest(t, CreateRandomUser(t))
	user := CreateRandomUser(t)
	for i := 0; i < 10; i++ {
		CreateRandomBuyRequest(t, sellRequest, user)
	}
	arg := ListBuyRequestsByTelegramIdParams{
		TelegramID: user.TelegramID,
		Limit: 5,
		Offset: 0,
	}
	buyRequests, err := testStore.ListBuyRequestsByTelegramId(context.Background(), arg)
	require.NoError(t, err)
	for _, buyRequest := range buyRequests {
		require.Equal(t, user.TelegramID, buyRequest.TelegramID)
	}
	require.Len(t, buyRequests, 5)
}

func TestOpenCloseBuyRequest(t *testing.T) {
	buyRequest1 := CreateRandomBuyRequest(t, CreateRandomSellRequest(t, CreateRandomUser(t)), CreateRandomUser(t))
	arg := ChangeStateBuyRequestParams{
		BuyReqID: buyRequest1.BuyReqID,
		State: "closed",
	}
	arg1 := CloseConfirmByBuyerParams{
		CloseConfirmByBuyer: util.ToPgBool(true),
		BuyReqID:            buyRequest1.BuyReqID,
	}
	arg2 := CloseConfirmBySellerParams{
		CloseConfirmBySeller: util.ToPgBool(true),
		BuyReqID:             buyRequest1.BuyReqID,
	}
	err := testStore.CloseConfirmByBuyer(context.Background(), arg1)
	require.NoError(t, err)
	err = testStore.CloseConfirmBySeller(context.Background(), arg2)
	require.NoError(t, err)
	buyRequest2, err := testStore.ChangeStateBuyRequest(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, buyRequest2.BuyReqID, buyRequest1.BuyReqID)
	require.Equal(t, "open", buyRequest1.State)
	require.Equal(t, "closed", buyRequest2.State)
}

func TestDeleteBuyRequest(t *testing.T) {
	buyRequest := CreateRandomBuyRequest(t, CreateRandomSellRequest(t, CreateRandomUser(t)), CreateRandomUser(t))
	err := testStore.DeleteBuyRequest(context.Background(), buyRequest.BuyReqID)
	require.NoError(t, err)
	_, err = testStore.GetBuyRequestById(context.Background(), buyRequest.BuyReqID)
	require.ErrorContains(t, sql.ErrNoRows, err.Error())
}
