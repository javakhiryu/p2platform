package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	mockdb "p2platform/db/mock"
	db "p2platform/db/sqlc"
	util "p2platform/util"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateSellRequest(t *testing.T) {
	sellRequest := randomSellRequest()
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"sell_total_amount":   sellRequest.SellTotalAmount,
				"currency_from":       sellRequest.CurrencyFrom,
				"currency_to":         sellRequest.CurrencyTo,
				"tg_username":         sellRequest.TgUsername,
				"sell_amount_by_card": sellRequest.SellAmountByCard,
				"sell_amount_by_cash": sellRequest.SellAmountByCash,
				"sell_exchange_rate":  sellRequest.SellExchangeRate,
				"comment":             sellRequest.Comment,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateSellRequestParams{
					SellTotalAmount:  sellRequest.SellTotalAmount,
					CurrencyFrom:     sellRequest.CurrencyFrom,
					CurrencyTo:       sellRequest.CurrencyTo,
					TgUsername:       sellRequest.TgUsername,
					SellByCard:       sellRequest.SellByCard,
					SellAmountByCard: sellRequest.SellAmountByCard,
					SellByCash:       sellRequest.SellByCash,
					SellAmountByCash: sellRequest.SellAmountByCash,
					SellExchangeRate: sellRequest.SellExchangeRate,
					Comment:          sellRequest.Comment,
				}
				store.EXPECT().
					CreateSellRequest(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(sellRequest, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchSellRequest(t, recorder.Body, sellRequest)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/createSellRequest")
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomSellRequest() db.SellRequest {
	return db.SellRequest{
		SellReqID:        int32(util.RandomInt(1, 1000)),
		SellTotalAmount:  1000,
		CurrencyFrom:     util.RandomCurrency(),
		CurrencyTo:       util.RandomCurrency(),
		TgUsername:       util.RandomTgUsername(),
		SellByCard:       util.ToPgBool(true),
		SellAmountByCard: util.ToPgInt(500),
		SellByCash:       util.ToPgBool(true),
		SellAmountByCash: util.ToPgInt(500),
		SellExchangeRate: util.ToPgInt(12950),
		Comment:          util.RandomString(10),
	}
}

func requireBodyMatchSellRequest(t *testing.T, body *bytes.Buffer, sellRequest db.SellRequest) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotSellRequest db.SellRequest
	err = json.Unmarshal(data, &gotSellRequest)
	require.NoError(t, err)
	require.Equal(t, sellRequest, gotSellRequest)
}
