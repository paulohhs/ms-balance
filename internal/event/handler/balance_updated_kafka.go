package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/paulohhs/ms-balance/internal/usecase/update_balance"
	"github.com/paulohhs/ms-balance/pkg/events"
)

type BalanceUpdatedKafkaHandler struct {
	UpdateBalanceUseCase *update_balance.UpdateBalanceUseCase
}

func NewBalanceUpdatedKafkaHandler(updateBalanceUseCase *update_balance.UpdateBalanceUseCase) *BalanceUpdatedKafkaHandler {
	return &BalanceUpdatedKafkaHandler{
		UpdateBalanceUseCase: updateBalanceUseCase,
	}
}

func (h *BalanceUpdatedKafkaHandler) Handle(message events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()

	payloadJSON, err := json.Marshal(message.GetPayload())
	if err != nil {
		fmt.Println("BalanceUpdatedKafkaHandler: invalid payload:", err)
		return
	}

	var input update_balance.UpdateBalanceInputDTO
	if err := json.Unmarshal(payloadJSON, &input); err != nil {
		fmt.Println("BalanceUpdatedKafkaHandler: invalid payload:", err)
		return
	}

	if err := h.UpdateBalanceUseCase.Execute(context.Background(), input); err != nil {
		fmt.Println("BalanceUpdatedKafkaHandler: error updating balance:", err)
		return
	}

	fmt.Println("BalanceUpdatedKafkaHandler", input)
}
