package api

import (
	"bytes"
	"database/sql"
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

func TestCreateSellRequestAPI(t *testing.T) {
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
				"sell_money_source":   sellRequest.SellMoneySource,
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
					SellMoneySource:  sellRequest.SellMoneySource,
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
		{
			name: "Bad Request",
			body: gin.H{
				"currency_from":       sellRequest.CurrencyFrom,
				"currency_to":         sellRequest.CurrencyTo,
				"tg_username":         sellRequest.TgUsername,
				"sell_amount_by_card": sellRequest.SellAmountByCard,
				"sell_amount_by_cash": sellRequest.SellAmountByCash,
				"sell_exchange_rate":  sellRequest.SellExchangeRate,
				"comment":             sellRequest.Comment,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateSellRequest(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Wrong total amount",
			body: gin.H{
				"sell_total_amount":   sellRequest.SellTotalAmount,
				"currency_from":       sellRequest.CurrencyFrom,
				"currency_to":         sellRequest.CurrencyTo,
				"tg_username":         sellRequest.TgUsername,
				"sell_amount_by_card": util.ToPgInt(100),
				"sell_amount_by_cash": sellRequest.SellAmountByCash,
				"sell_exchange_rate":  sellRequest.SellExchangeRate,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateSellRequest(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.NotEmpty(t, recorder.Body)
			},
		},
		{
			name: "Internal Server Error",
			body: gin.H{
				"sell_total_amount":   sellRequest.SellTotalAmount,
				"sell_money_source":   sellRequest.SellMoneySource,
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
					SellMoneySource:  sellRequest.SellMoneySource,
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
					Return(db.SellRequest{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			url := fmt.Sprintf("/sell-request")
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetSellRequestAPI(t *testing.T) {
	sellRequest := randomSellRequest()
	deletedSellRequest := deletedSellRequest()
	testCases := []struct {
		name          string
		sellRequestID int32
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:          "OK",
			sellRequestID: sellRequest.SellReqID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSellRequestById(gomock.Any(), gomock.Eq(sellRequest.SellReqID)).
					Times(1).
					Return(sellRequest, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:          "Bad Request",
			sellRequestID: -1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetSellRequestById(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:          "Internal Server Error",
			sellRequestID: sellRequest.SellReqID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetSellRequestById(gomock.Any(), gomock.Eq(sellRequest.SellReqID)).Times(1).Return(db.SellRequest{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:          "Deleted Sell Request",
			sellRequestID: deletedSellRequest.SellReqID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetSellRequestById(gomock.Any(), gomock.Eq(deletedSellRequest.SellReqID)).Times(1).Return(deletedSellRequest, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusGone, recorder.Code)
				require.NotEmpty(t, recorder.Body)
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

			url := fmt.Sprintf("/sell-request/%d", tc.sellRequestID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestListSellRequestsAPI(t *testing.T) {

	n := 10
	var sellRequests = make([]db.SellRequest, n)
	for i := 0; i < n; i++ {
		sellRequests[i] = randomSellRequest()
	}

	type Query struct {
		PageSize int
		PageId   int
	}

	testCases := []struct {
		name          string
		query         Query
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				PageSize: n,
				PageId:   1,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListSellRequestsParams{
					Limit:  int32(n),
					Offset: 0,
				}
				store.EXPECT().ListSellRequests(gomock.Any(), gomock.Eq(arg)).Times(1).Return(sellRequests, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchSellRequests(t, recorder.Body, sellRequests)
			},
		},
		{
			name: "Bad Request",
			query: Query{
				PageSize: 20,
				PageId:   0,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().ListSellRequests(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Internal Server Error",
			query: Query{
				PageSize: n,
				PageId:   1,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListSellRequestsParams{
					Limit:  int32(n),
					Offset: 0,
				}
				store.EXPECT().ListSellRequests(gomock.Any(), gomock.Eq(arg)).Times(1).Return(sellRequests, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			url := fmt.Sprintf("/sell-requests")
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.PageId))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.PageSize))
			request.URL.RawQuery = q.Encode()
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)

		})
	}
}

func TestUpdateSellRequestAPI(t *testing.T) {
	sellRequest := randomSellRequest()
	deletedSellRequest := deletedSellRequest()
	testCases := []struct {
		name          string
		sellRequestID int32
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:          "OK",
			sellRequestID: sellRequest.SellReqID,
			body: gin.H{
				"sell_exchange_rate": 1,
				"comment":            "new comment",
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateSellRequestParams{
					SellReqID:        sellRequest.SellReqID,
					SellExchangeRate: util.ToPgInt(1),
					Comment:          util.ToPgText("new comment"),
					SellByCard:       util.ToPgBool(true),
					SellByCash:       util.ToPgBool(true),
				}

				store.EXPECT().GetSellRequestById(gomock.Any(), gomock.Eq(sellRequest.SellReqID)).Times(1).Return(sellRequest, nil)
				store.EXPECT().UpdateSellRequest(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.SellRequest{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:          "Incorrect ID",
			sellRequestID: -1,
			body: gin.H{
				"sell_exchange_rate": 1,
				"comment":            "new comment",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetSellRequestById(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().UpdateSellRequest(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:          "Incorrect currency",
			sellRequestID: sellRequest.SellReqID,
			body: gin.H{
				"currency_from": "test",
				"comment":       "new comment",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetSellRequestById(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().UpdateSellRequest(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:          "Internal Server Error on Get",
			sellRequestID: sellRequest.SellReqID,
			body: gin.H{
				"sell_exchange_rate": 1,
				"comment":            "new comment",
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateSellRequestParams{
					SellReqID:        sellRequest.SellReqID,
					SellExchangeRate: util.ToPgInt(1),
					Comment:          util.ToPgText("new comment"),
				}
				store.EXPECT().GetSellRequestById(gomock.Any(), gomock.Eq(sellRequest.SellReqID)).Times(1).Return(db.SellRequest{}, sql.ErrConnDone)
				store.EXPECT().UpdateSellRequest(gomock.Any(), gomock.Eq(arg)).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:          "Update deleted sell request",
			sellRequestID: deletedSellRequest.SellReqID,
			body: gin.H{
				"sell_exchange_rate": 1,
				"comment":            "new comment",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetSellRequestById(gomock.Any(), gomock.Eq(deletedSellRequest.SellReqID)).Times(1).Return(deletedSellRequest, nil)
				store.EXPECT().UpdateSellRequest(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusGone, recorder.Code)
			},
		},
		{
			name:          "Wrong total amount",
			sellRequestID: sellRequest.SellReqID,
			body: gin.H{
				"sell_total_amount":   1000,
				"sell_amount_by_card": 1000,
				"sell_amount_by_cash": 100,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetSellRequestById(gomock.Any(), gomock.Eq(sellRequest.SellReqID)).Times(1).Return(sellRequest, nil)
				store.EXPECT().UpdateSellRequest(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:          "Sell by card",
			sellRequestID: sellRequest.SellReqID,
			body: gin.H{
				"sell_total_amount":   900,
				"sell_amount_by_card": 900,
				"sell_amount_by_cash": 0,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateSellRequestParams{
					SellReqID:        sellRequest.SellReqID,
					SellTotalAmount:  util.ToPgInt(900),
					SellAmountByCard: util.ToPgInt(900),
					SellAmountByCash: util.ToPgInt(0),
					SellByCard:       util.ToPgBool(true),
					SellByCash:       util.ToPgBool(false),
				}
				store.EXPECT().GetSellRequestById(gomock.Any(), gomock.Eq(sellRequest.SellReqID)).Times(1).Return(sellRequest, nil)
				store.EXPECT().UpdateSellRequest(gomock.Any(), gomock.Eq(arg)).
					Times(1).Return(sellRequest, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:          "Sell by cash",
			sellRequestID: sellRequest.SellReqID,
			body: gin.H{
				"sell_total_amount":   900,
				"sell_amount_by_card": 0,
				"sell_amount_by_cash": 900,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateSellRequestParams{
					SellReqID:        sellRequest.SellReqID,
					SellTotalAmount:  util.ToPgInt(900),
					SellAmountByCard: util.ToPgInt(0),
					SellAmountByCash: util.ToPgInt(900),
					SellByCard:       util.ToPgBool(false),
					SellByCash:       util.ToPgBool(true),
				}
				store.EXPECT().GetSellRequestById(gomock.Any(), gomock.Eq(sellRequest.SellReqID)).
					Times(1).
					Return(sellRequest, nil)
				store.EXPECT().UpdateSellRequest(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(sellRequest, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:          "Internal Server Error on Update",
			sellRequestID: sellRequest.SellReqID,
			body: gin.H{
				"sell_exchange_rate": 1,
				"comment":            "new comment",
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateSellRequestParams{
					SellReqID:        sellRequest.SellReqID,
					SellExchangeRate: util.ToPgInt(1),
					Comment:          util.ToPgText("new comment"),
					SellByCard:       util.ToPgBool(true),
					SellByCash:       util.ToPgBool(true),
				}
				store.EXPECT().GetSellRequestById(gomock.Any(), gomock.Eq(sellRequest.SellReqID)).Times(1).Return(sellRequest, nil)
				store.EXPECT().UpdateSellRequest(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.SellRequest{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			url := fmt.Sprintf("/sell-request/%d", tc.sellRequestID)
			request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestDeleteSellRequest(t *testing.T) {
	sellRequest := randomSellRequest()
	deletedSellRequest := deletedSellRequest()
	testCases := []struct {
		name          string
		sellReqID     int32
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			sellReqID: sellRequest.SellReqID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteSellRequestTx(gomock.Any(), gomock.Eq(sellRequest.SellReqID)).Times(1).Return(true, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "Bad Request",
			sellReqID: -1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteSellRequestTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "Internal Server Error On Update",
			sellReqID: sellRequest.SellReqID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteSellRequestTx(gomock.Any(), gomock.Any()).Times(1).Return(false, sql.ErrTxDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "Deleted Sell Request",
			sellReqID: deletedSellRequest.SellReqID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteSellRequestTx(gomock.Any(), gomock.Eq(deletedSellRequest.SellReqID)).Times(1).Return(false, db.ErrSellRequestAlreadyDeleted)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, recorder.Code)
			},
		},
		{
			name:      "Sell Request Not Found",
			sellReqID: 1000,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteSellRequestTx(gomock.Any(), gomock.Eq(int32(1000))).Times(1).Return(false, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
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

			url := fmt.Sprintf("/sell-request/%d", tc.sellReqID)

			server := newTestServer(t, store)

			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomSellRequest() db.SellRequest {
	return db.SellRequest{
		SellReqID:        int32(util.RandomInt(1, 1000)),
		SellTotalAmount:  1000,
		SellMoneySource:  "cash",
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

func deletedSellRequest() db.SellRequest {
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
		IsDeleted:        util.ToPgBool(true),
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

func requireBodyMatchSellRequests(t *testing.T, body *bytes.Buffer, sellRequests []db.SellRequest) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotSellRequests []db.SellRequest
	err = json.Unmarshal(data, &gotSellRequests)
	require.NoError(t, err)
	require.Equal(t, sellRequests, gotSellRequests)
}
