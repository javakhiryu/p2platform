package db

import (
	"context"
	"database/sql"
	"p2platform/util"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomBuyRequest(t *testing.T) BuyRequest { 
	sellRequest := createRandomSellRequest(t)
	arg := CreateBuyRequestParams{
		BuyReqID: uuid.New(),
		SellReqID: sellRequest.SellReqID,
		BuyAmount: sellRequest.SellAmount,
		TgUsername: util.RandomTgUsername(),
		BuyByCard: sellRequest.SellByCard,
		BuyAmountByCard: sellRequest.SellAmountByCard,
		BuyByCash: sellRequest.SellByCash,
		BuyAmountByCash: sellRequest.SellAmountByCash,

	}
	buyRequest, err := testStore.CreateBuyRequest(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, buyRequest)
	require.Equal(t, sellRequest.SellReqID,  buyRequest.SellReqID)
	return buyRequest
}

func TestCreateBuyRequest(t *testing.T){
	createRandomBuyRequest(t)
}

func TestGetBuyRequest(t *testing.T){
	buyRequest := createRandomBuyRequest(t)
	buyRequest2, err := testStore.GetBuyRequestById(context.Background(), buyRequest.BuyReqID)
	require.NoError(t, err)
	require.NotEmpty(t, buyRequest2)
	require.Equal(t, buyRequest.BuyReqID, buyRequest2.BuyReqID)
}

func TestListBuyRequests(t *testing.T){
	var lastBuyRequest BuyRequest
	for i:=0; i<5; i++{
		lastBuyRequest = createRandomBuyRequest(t)
	}
	arg := ListBuyRequestsParams{
		SellReqID: lastBuyRequest.SellReqID,
		Limit: 5,
		Offset: 0,
	}
	buyRequests, err := testStore.ListBuyRequests(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, buyRequests)
	for _,br :=range buyRequests{
		require.Equal(t, lastBuyRequest.BuyReqID, br.BuyReqID)
	}
} 

func TestUpdateBuyRequest(t *testing.T){
	buyRequest1 := createRandomBuyRequest(t)
	newTgUsername := util.RandomTgUsername()
	arg := UpdateBuyRequestParams{
		BuyReqID: buyRequest1.BuyReqID,
		TgUsername: newTgUsername,
	}
	buyRequest2, err := testStore.UpdateBuyRequest(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, buyRequest2.BuyReqID, buyRequest1.BuyReqID)
	require.NotEqual(t, buyRequest2.TgUsername, buyRequest1.TgUsername)
	require.Equal(t, newTgUsername, buyRequest2.TgUsername)
}

func TestOpenCloseBuyRequest(t *testing.T){
	buyRequest1 := createRandomBuyRequest(t)
	arg := OpenCloseBuyRequestParams{
		BuyReqID: buyRequest1.BuyReqID,
		IsSuccessful: util.ToPgBool(true),
	}
	buyRequest2, err := testStore.OpenCloseBuyRequest(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, buyRequest2.BuyReqID, buyRequest1.BuyReqID)
	require.Equal(t, util.ToPgBool(false), buyRequest1.IsSuccessful)
	require.Equal(t, util.ToPgBool(true), buyRequest2.IsSuccessful) 
}

func TestDeleteBuyRequest(t *testing.T){
	buyRequest := createRandomBuyRequest(t)
	err :=testStore.DeleteBuyRequest(context.Background(), buyRequest.BuyReqID)
	require.NoError(t, err)
	_, err = testStore.GetBuyRequestById(context.Background(), buyRequest.BuyReqID)
	require.ErrorContains(t, sql.ErrNoRows, err.Error())
}


