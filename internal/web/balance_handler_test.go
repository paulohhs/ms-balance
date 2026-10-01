package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/paulohhs/ms-balance/internal/entity"
	"github.com/paulohhs/ms-balance/internal/usecase/get_balance"
	"github.com/paulohhs/ms-balance/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
)

func serve(h *WebBalanceHandler, path string) *httptest.ResponseRecorder {
	router := chi.NewRouter()
	router.Get("/balances/{account_id}", h.GetBalance)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec
}

func TestWebBalanceHandler_GetBalance(t *testing.T) {
	account, _ := entity.NewAccount("a1", 1000)
	accountMock := &mocks.AccountGatewayMock{}
	accountMock.On("FindByID", "a1").Return(account, nil)
	h := NewWebBalanceHandler(*get_balance.NewGetBalanceUseCase(accountMock))

	rec := serve(h, "/balances/a1")

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	var body get_balance.GetBalanceOutputDTO
	assert.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
	assert.Equal(t, "a1", body.AccountID)
	assert.Equal(t, float64(1000), body.Balance)
}

func TestWebBalanceHandler_GetBalance_NotFound(t *testing.T) {
	accountMock := &mocks.AccountGatewayMock{}
	accountMock.On("FindByID", "x").Return(nil, entity.ErrAccountNotFound)
	h := NewWebBalanceHandler(*get_balance.NewGetBalanceUseCase(accountMock))

	rec := serve(h, "/balances/x")

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestWebBalanceHandler_GetBalance_InternalError(t *testing.T) {
	accountMock := &mocks.AccountGatewayMock{}
	accountMock.On("FindByID", "a1").Return(nil, errors.New("db down"))
	h := NewWebBalanceHandler(*get_balance.NewGetBalanceUseCase(accountMock))

	rec := serve(h, "/balances/a1")

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
