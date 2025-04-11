package db

import (
	"context"
	"p2platform/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomSellRequest(t *testing.T) SellRequest {
	SellAmount := util.RandomMoney()
	arg := CreateSellRequestParams{
		SellTotalAmount:  SellAmount,
		SellMoneySource:  "card",
		CurrencyFrom:     util.RandomCurrency(),
		CurrencyTo:       util.RandomCurrency(),
		TgUsername:       util.RandomTgUsername(),
		SellAmountByCard: util.ToPgInt(SellAmount / 2),
		SellAmountByCash: util.ToPgInt(SellAmount / 2),
		SellExchangeRate: util.ToPgInt(12950),
		Comment:          util.RandomString(10),
	}
	sellRequest, err := testStore.CreateSellRequest(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, sellRequest)
	require.Equal(t, arg.SellTotalAmount, sellRequest.SellTotalAmount)
	require.Equal(t, arg.SellMoneySource, sellRequest.SellMoneySource)
	require.Equal(t, arg.CurrencyFrom, sellRequest.CurrencyFrom)
	require.Equal(t, arg.CurrencyTo, sellRequest.CurrencyTo)
	require.Equal(t, arg.TgUsername, sellRequest.TgUsername)
	require.Equal(t, arg.SellAmountByCard, sellRequest.SellAmountByCard)
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
	require.Equal(t, sellRequest1.SellTotalAmount, sellRequest2.SellTotalAmount)
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
	newSellTotalAmount := util.RandomMoney()
	newCurrencyFrom := util.RandomCurrency()
	newCurrencyTo := util.RandomCurrency()
	newComment := util.RandomString(10)

	sellRequest1 := createRandomSellRequest(t)
	arg := UpdateSellRequestParams{
		SellReqID:       sellRequest1.SellReqID,
		SellTotalAmount: util.ToPgInt(newSellTotalAmount),
		CurrencyFrom:    util.ToPgText(newCurrencyFrom),
		CurrencyTo:      util.ToPgText(newCurrencyTo),
		Comment:         util.ToPgText(newComment),
	}
	sellRequest2, err := testStore.UpdateSellRequest(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, sellRequest2)
	require.Equal(t, arg.SellReqID, sellRequest2.SellReqID)
	require.Equal(t, sellRequest1.TgUsername, sellRequest2.TgUsername)
	require.Equal(t, newSellTotalAmount, sellRequest2.SellTotalAmount)
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

func TestCloseSellRequest(t *testing.T) {
	sellRequest1 := createRandomSellRequest(t)
	arg := OpenCloseSellRequestParams{
		IsActual:  util.ToPgBool(false),
		SellReqID: sellRequest1.SellReqID,
	}
	sellRequest2, err := testStore.OpenCloseSellRequest(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, sellRequest2)
	require.Equal(t, util.ToPgBool(true), sellRequest1.IsActual)
	require.Equal(t, util.ToPgBool(false), sellRequest2.IsActual)
	require.NotEqual(t, sellRequest2.UpdatedAt, sellRequest1.UpdatedAt)
}

func TestDeleteSellRequest(t *testing.T) {
	SellRequest1 := createRandomSellRequest(t)
	isDeleted, err := testStore.DeleteSellRequest(context.Background(), SellRequest1.SellReqID)
	require.NoError(t, err)
	require.Equal(t, util.ToPgBool(false), SellRequest1.IsDeleted)
	require.Equal(t, util.ToPgBool(true), isDeleted)
}

func TestGetSellRequestForUpdate_Lock(t *testing.T) {
	sellRequest := createRandomSellRequest(t)
	sqlStore := testStore.(*SQLStore)
	conn1, err := sqlStore.connPool.Acquire(context.Background())
	require.NoError(t, err)
	defer conn1.Release()

	tx1, err := conn1.Begin(context.Background())
	require.NoError(t, err)
	q1 := New(tx1)

	_, err = q1.GetSellRequestForUpdate(context.Background(), sellRequest.SellReqID)
	require.NoError(t, err)

	done := make(chan struct{})
	go func() {
		conn2, err := sqlStore.connPool.Acquire(context.Background())
		require.NoError(t, err)
		defer conn2.Release()

		tx2, err := conn2.Begin(context.Background())
		require.NoError(t, err)
		q2 := New(tx2)

		start := time.Now()
		_, err = q2.GetSellRequestForUpdate(context.Background(), sellRequest.SellReqID)
		require.NoError(t, err)
		duration := time.Since(start)

		require.GreaterOrEqual(t, duration.Milliseconds(), int64(70))

		tx2.Commit(context.Background())
		close(done)
	}()

	time.Sleep(100 * time.Millisecond)

	err = tx1.Commit(context.Background())
	require.NoError(t, err)

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("second transaction did not finish in time")
	}
}
