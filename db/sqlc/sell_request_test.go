package db

import (
	"context"
	"p2platform/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomSellRequest(t *testing.T) SellRequest {
	SellAmount := util.RandomMoney()
	arg := CreateSellRequestParams{
		SellAmount:       SellAmount,
		CurrencyFrom:     util.RandomCurrency(),
		CurrencyTo:       util.RandomCurrency(),
		TgUsername:       util.RandomTgUsername(),
		SellByCard:       util.ToPgBool(true),
		SellAmountByCard: util.ToPgInt(SellAmount/2),
		SellByCash:       util.ToPgBool(true),
		SellAmountByCash:  util.ToPgInt(SellAmount/2),
		SellExchangeRate: util.ToPgInt(12950),
		Comment:          util.RandomString(10),
	}
	sellRequest, err := testStore.CreateSellRequest(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, sellRequest)
	require.Equal(t, arg.SellAmount, sellRequest.SellAmount)
	require.Equal(t, arg.CurrencyFrom, sellRequest.CurrencyFrom)
	require.Equal(t, arg.CurrencyTo, sellRequest.CurrencyTo)
	require.Equal(t, arg.TgUsername, sellRequest.TgUsername)
	require.Equal(t, arg.SellAmountByCard, sellRequest.SellAmountByCard)
	require.Equal(t, arg.SellByCash, sellRequest.SellByCash)
	require.Equal(t, arg.SellAmountByCash, sellRequest.SellAmountByCash)
	require.Equal(t, arg.SellExchangeRate, sellRequest.SellExchangeRate)
	require.Equal(t, arg.Comment, sellRequest.Comment)

	return sellRequest

}

func TestCreateSellReqest(t *testing.T) {
	createRandomSellRequest(t)
}

func TestGetSellRequest(t *testing.T) {
	sellRequest1 := createRandomSellRequest(t)
	sellRequest2, err := testStore.GetSellRequestById(context.Background(), sellRequest1.SellReqID)
	require.NoError(t, err)
	require.NotEmpty(t, sellRequest2)
	require.Equal(t, sellRequest1.SellReqID, sellRequest2.SellReqID)
	require.Equal(t, sellRequest1.SellAmount, sellRequest2.SellAmount)
	require.Equal(t, sellRequest1.CurrencyFrom, sellRequest2.CurrencyFrom)
	require.Equal(t, sellRequest1.CurrencyTo, sellRequest2.CurrencyTo)
	require.Equal(t, sellRequest1.TgUsername, sellRequest2.TgUsername)
	require.Equal(t, sellRequest1.SellByCard, sellRequest2.SellByCard)
	require.Equal(t, sellRequest1.SellAmountByCard, sellRequest2.SellAmountByCard)
	require.Equal(t, sellRequest1.SellByCash, sellRequest2.SellByCash)
	require.Equal(t, sellRequest1.SellAmountByCash, sellRequest2.SellAmountByCash)
	require.Equal(t, sellRequest1.SellExchangeRate, sellRequest2.SellExchangeRate)
	require.Equal(t, sellRequest1.Comment, sellRequest2.Comment)

}

func TestListSellRequests(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomSellRequest(t)
	}
	arg := ListSellRequestsParams{
		Limit:  5,
		Offset: 0,
	}
	sellRequests, err := testStore.ListSellRequests(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, sellRequests, int(arg.Limit))
	for _, sellRequest := range sellRequests {
		require.NotEmpty(t, sellRequest)
	}
}

func TestUpdateSellRequest(t *testing.T) {
	newSellAmount := util.RandomMoney()
	newCurrencyFrom :=util.RandomCurrency()
	newCurrencyTo := util.RandomCurrency()
	newComment := util.RandomString(10)

	sellRequest1 := createRandomSellRequest(t)
	arg := UpdateSellRequestParams{
		SellReqID:        sellRequest1.SellReqID,
		SellAmount:       util.ToPgInt(newSellAmount),
		CurrencyFrom:     util.ToPgText(newCurrencyFrom),
		CurrencyTo:       util.ToPgText(newCurrencyTo),
		Comment: 		  util.ToPgText(newComment),

		
	}
	sellRequest2, err := testStore.UpdateSellRequest(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, sellRequest2)
	require.Equal(t, arg.SellReqID, sellRequest2.SellReqID)
	require.Equal(t, sellRequest1.TgUsername, sellRequest2.TgUsername)
	require.Equal(t, newSellAmount, sellRequest2.SellAmount)
	require.Equal(t, newCurrencyFrom, sellRequest2.CurrencyFrom)
	require.Equal(t, newCurrencyTo, sellRequest2.CurrencyTo)
	require.Equal(t, sellRequest1.SellByCard, sellRequest2.SellByCard)
	require.Equal(t, sellRequest1.SellAmountByCard, sellRequest2.SellAmountByCard)
	require.Equal(t, sellRequest1.SellByCash, sellRequest2.SellByCash)
	require.Equal(t, sellRequest1.SellAmountByCash, sellRequest2.SellAmountByCash)
	require.Equal(t, sellRequest1.SellExchangeRate, sellRequest2.SellExchangeRate)
	require.Equal(t, newComment, sellRequest2.Comment)
	require.NotEqual(t, sellRequest1.Comment, sellRequest2.Comment)
	require.NotEqual(t, sellRequest2.CreatedAt, sellRequest2.UpdatedAt)

}

func TestCloseSellRequest(t *testing.T){
	sellRequest1 := createRandomSellRequest(t)
	arg := OpenCloseSellRequestParams{
		IsActual: util.ToPgBool(false),
		SellReqID: sellRequest1.SellReqID,
	}
	sellRequest2, err := testStore.OpenCloseSellRequest(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, sellRequest2)
	require.Equal(t, util.ToPgBool(true), sellRequest1.IsActual)
	require.Equal(t, util.ToPgBool(false), sellRequest2.IsActual)
	require.NotEqual(t, sellRequest2.UpdatedAt, sellRequest1.UpdatedAt)
}

func TestDeleteSellRequest(t *testing.T){
	SellRequest1 := createRandomSellRequest(t)
	SellRequest2, err := testStore.DeleteSellRequest(context.Background(), SellRequest1.SellReqID)
	require.NoError(t, err)
	require.Equal(t, util.ToPgBool(false), SellRequest1.IsDeleted)
	require.Equal(t, util.ToPgBool(true), SellRequest2.IsDeleted)
}