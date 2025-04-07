package db

import (
	"context"
	"p2platform/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomLockedAmount(t *testing.T, buyRequest BuyRequest) LockedAmount {
	arg := CreateLockedAmountParams{
		SellReqID:         buyRequest.SellReqID,
		BuyReqID:          buyRequest.BuyReqID,
		LockedTotalAmount: buyRequest.BuyTotalAmount,
		LockedByCard:      buyRequest.BuyAmountByCard,
		LockedByCash:      buyRequest.BuyAmountByCash,
	}
	lockedAmount, err := testStore.CreateLockedAmount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, lockedAmount)

	return lockedAmount
}

func TestCreateLockedAmount(t *testing.T) {
	createRandomLockedAmount(t, createRandomBuyRequest(t, createRandomSellRequest(t)))
}

func TestGetLockedAmount(t *testing.T) {
	lockedAmount := createRandomLockedAmount(t, createRandomBuyRequest(t, createRandomSellRequest(t)))
	lockedAmount2, err := testStore.GetLockedAmount(context.Background(), lockedAmount.BuyReqID)
	require.NoError(t, err)
	require.NotEmpty(t, lockedAmount2)
	require.Equal(t, lockedAmount.BuyReqID, lockedAmount2.BuyReqID)
}
func TestGetLockedAmountBySellRequest(t *testing.T) {
	sellRequest := createRandomSellRequest(t)
	lockedAmount := createRandomLockedAmount(t, createRandomBuyRequest(t, sellRequest))
	lockedAmount2, err := testStore.GetLockedAmountBySellRequest(context.Background(), sellRequest.SellReqID)
	require.NoError(t, err)
	require.NotEmpty(t, lockedAmount2)
	for _,la :=range lockedAmount2{
		require.Equal(t, sellRequest.SellReqID, la.SellReqID)
		require.Equal(t, lockedAmount.ID, la.ID)
	}
}

func TestListLockedAmounts(t *testing.T) {
	sellRequest := createRandomSellRequest(t)
	var lockedAmounts []LockedAmount
	for i:=0; i<5; i++{
		lockedAmounts = append(lockedAmounts, createRandomLockedAmount(t, createRandomBuyRequest(t, sellRequest)))
	}
	arg := ListLockedAmountsParams{
		SellReqID: lockedAmounts[0].SellReqID,
		Limit: 5,
		Offset: 0,
	}
	lockedAmounts, err := testStore.ListLockedAmounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, lockedAmounts)
	for _,la :=range lockedAmounts{
		require.Equal(t, sellRequest.SellReqID, la.SellReqID)
	}
}

func TestReleaseLockedAmountByBuyRequest(t *testing.T){
	sellRequest := createRandomSellRequest(t)
	buyRequest := createRandomBuyRequest(t, sellRequest) 
	lockedAmount1 := createRandomLockedAmount(t, buyRequest)
	err := testStore.ReleaseLockedAmountByBuyRequest(context.Background(), buyRequest.BuyReqID)
	require.NoError(t, err)
	require.Equal(t, util.ToPgBool(false), lockedAmount1.IsReleased)
	lockedAmount2, err := testStore.GetLockedAmount(context.Background(), buyRequest.BuyReqID)
	require.NoError(t, err)
	require.Equal(t, util.ToPgBool(true), lockedAmount2.IsReleased)

}

func TestReleaseLockedAmountBySellRequest(t *testing.T){
	sellRequest := createRandomSellRequest(t)
	var buyRequests []BuyRequest
	for i:=0; i<5; i++{
		buyRequests = append(buyRequests, createRandomBuyRequest(t, sellRequest))
	}
	var lockedAmounts []LockedAmount
	for i:=0; i<5; i++{
		lockedAmounts = append(lockedAmounts, createRandomLockedAmount(t, buyRequests[i]))
	}

	for _,la :=range lockedAmounts{
		require.Equal(t, util.ToPgBool(false), la.IsReleased)
	}

	err := testStore.ReleaseLockedAmountsBySellRequest(context.Background(), sellRequest.SellReqID)
	require.NoError(t, err)

	lockedAmounts2, err := testStore.GetLockedAmountBySellRequest(context.Background(), sellRequest.SellReqID)
	require.NoError(t, err)
	for _,la :=range lockedAmounts2{
		require.Equal(t, util.ToPgBool(true), la.IsReleased)
	}
	
	
}