package handler

import (
	"encoding/json"
	"sync"
	"testing"

	"github.com/paulohhs/ms-balance/internal/entity"
	"github.com/paulohhs/ms-balance/internal/event"
	"github.com/paulohhs/ms-balance/internal/usecase/mocks"
	"github.com/paulohhs/ms-balance/internal/usecase/update_balance"
	"github.com/paulohhs/ms-balance/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBalanceUpdatedKafkaHandler_Handle(t *testing.T) {
	accountFrom, _ := entity.NewAccount("from", 1000)
	accountTo, _ := entity.NewAccount("to", 1000)
	accountMock := &mocks.AccountGatewayMock{}
	accountMock.On("FindByID", "from").Return(accountFrom, nil)
	accountMock.On("FindByID", "to").Return(accountTo, nil)
	accountMock.On("UpdateBalance", mock.Anything).Return(nil)

	mockUow := &mocks.UowMock{}
	mockUow.On("Do", mock.Anything, mock.Anything).Return(nil)
	mockUow.On("GetRepository", mock.Anything, "AccountDB").Return(accountMock, nil)

	// Mensagem exatamente como o wallet publica no tópico "balances".
	balanceUpdated := event.NewBalanceUpdated()
	err := json.Unmarshal([]byte(`{"Name":"BalanceUpdated","Payload":{"account_id_from":"from","account_id_to":"to","balance_account_id_from":900,"balance_account_id_to":1100}}`), balanceUpdated)
	assert.NoError(t, err)

	// Passa pelo dispatcher real: prova que o handler está registrado no nome certo.
	dispatcher := events.NewEventDispatcher()
	dispatcher.Register("BalanceUpdated", NewBalanceUpdatedKafkaHandler(update_balance.NewUpdateBalanceUseCase(mockUow)))
	dispatcher.Dispatch(balanceUpdated)

	assert.Equal(t, float64(900), accountFrom.Balance)
	assert.Equal(t, float64(1100), accountTo.Balance)
	accountMock.AssertNumberOfCalls(t, "UpdateBalance", 2)
}

func TestBalanceUpdatedKafkaHandler_Handle_InvalidPayload(t *testing.T) {
	mockUow := &mocks.UowMock{}
	h := NewBalanceUpdatedKafkaHandler(update_balance.NewUpdateBalanceUseCase(mockUow))

	balanceUpdated := event.NewBalanceUpdated()
	balanceUpdated.SetPayload("payload-invalido")

	wg := &sync.WaitGroup{}
	wg.Add(1)
	h.Handle(balanceUpdated, wg)
	wg.Wait() // se o handler esquecer o wg.Done(), o teste trava até o timeout

	mockUow.AssertNotCalled(t, "Do", mock.Anything, mock.Anything)
}
