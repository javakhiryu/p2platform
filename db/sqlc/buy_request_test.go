package db

import (
	"context"
	"database/sql"
	"p2platform/util"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomBuyRequest(t *testing.T, sellRequest SellRequest) BuyRequest { 
	arg := CreateBuyRequestParams{
		BuyReqID: uuid.New(),
		SellReqID: sellRequest.SellReqID,
		BuyTotalAmount: sellRequest.SellTotalAmount,
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
	createRandomBuyRequest(t, createRandomSellRequest(t))
}

func TestGetBuyRequest(t *testing.T){
	buyRequest := createRandomBuyRequest(t, createRandomSellRequest(t))
	buyRequest2, err := testStore.GetBuyRequestById(context.Background(), buyRequest.BuyReqID)
	require.NoError(t, err)
	require.NotEmpty(t, buyRequest2)
	require.Equal(t, buyRequest.BuyReqID, buyRequest2.BuyReqID)
}

func TestListBuyRequests(t *testing.T){
	sellRequest := createRandomSellRequest(t) 
	var buyRequests []BuyRequest
	for i:=0; i<5; i++{
		buyRequests = append(buyRequests, createRandomBuyRequest(t, sellRequest)) 
	}
	arg := ListBuyRequestsParams{
		SellReqID: sellRequest.SellReqID,
		Limit: 5,
		Offset: 0,
	}
	buyRequests, err := testStore.ListBuyRequests(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, buyRequests)
	require.Len(t, buyRequests, 5)
	for _,br :=range buyRequests{
		require.Equal(t, buyRequests[0].SellReqID, br.SellReqID)
	}
} 

func TestOpenCloseBuyRequest(t *testing.T){
	buyRequest1 := createRandomBuyRequest(t, createRandomSellRequest(t))
	arg := OpenCloseBuyRequestParams{
		BuyReqID: buyRequest1.BuyReqID,
		IsClosed: util.ToPgBool(true),
	}
	arg1 := CloseConfirmByBuyerParams{
		CloseConfirmByBuyer: util.ToPgBool(true),
		BuyReqID: buyRequest1.BuyReqID,
	}
	arg2 := CloseConfirmBySellerParams{
		CloseConfirmBySeller: util.ToPgBool(true),
		BuyReqID: buyRequest1.BuyReqID,
	}
	err := testStore.CloseConfirmByBuyer(context.Background(), arg1)
	require.NoError(t, err)
	err = testStore.CloseConfirmBySeller(context.Background(), arg2)
	require.NoError(t, err)
	buyRequest2, err := testStore.OpenCloseBuyRequest(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, buyRequest2.BuyReqID, buyRequest1.BuyReqID)
	require.Equal(t, util.ToPgBool(false), buyRequest1.IsClosed)
	require.Equal(t, util.ToPgBool(true), buyRequest2.IsClosed) 
}

func TestDeleteBuyRequest(t *testing.T){
	buyRequest := createRandomBuyRequest(t, createRandomSellRequest(t))
	err :=testStore.DeleteBuyRequest(context.Background(), buyRequest.BuyReqID)
	require.NoError(t, err)
	_, err = testStore.GetBuyRequestById(context.Background(), buyRequest.BuyReqID)
	require.ErrorContains(t, sql.ErrNoRows, err.Error())
}


