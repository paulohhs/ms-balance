package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/paulohhs/ms-balance/internal/entity"
	"github.com/paulohhs/ms-balance/internal/usecase/get_balance"
)

type WebBalanceHandler struct {
	GetBalanceUseCase get_balance.GetBalanceUseCase
}

func NewWebBalanceHandler(getBalanceUseCase get_balance.GetBalanceUseCase) *WebBalanceHandler {
	return &WebBalanceHandler{
		GetBalanceUseCase: getBalanceUseCase,
	}
}

func (h *WebBalanceHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "account_id")

	output, err := h.GetBalanceUseCase.Execute(get_balance.GetBalanceInputDTO{AccountID: accountID})
	if errors.Is(err, entity.ErrAccountNotFound) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
