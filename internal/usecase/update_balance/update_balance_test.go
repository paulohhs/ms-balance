package updatebalance

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/paulohhs/ms-balance/internal/entity"
	"github.com/paulohhs/ms-balance/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateBalanceUseCase_Execute_UpdatesExistingAccounts(t *testing.T) {
	accountFrom, _ := entity.NewAccount("from", 1000)
	accountTo, _ := entity.NewAccount("to", 1000)

	accountMock := &mocks.AccountGatewayMock{}
	accountMock.On("FindByID", "from").Return(accountFrom, nil)
	accountMock.On("FindByID", "to").Return(accountTo, nil)
	accountMock.On("UpdateBalance", mock.Anything).Return(nil)

	uc := NewUpdateBalanceUseCase(accountMock)
	err := uc.Execute(UpdateBalanceInputDTO{
		AccountIDFrom:        "from",
		AccountIDTo:          "to",
		BalanceAccountIDFrom: 900,
		BalanceAccountIDTo:   1100,
	})

	assert.Nil(t, err)
	assert.Equal(t, float64(900), accountFrom.Balance)
	assert.Equal(t, float64(1100), accountTo.Balance)
	accountMock.AssertNumberOfCalls(t, "UpdateBalance", 2)
	accountMock.AssertNotCalled(t, "Save", mock.Anything)
}

func TestUpdateBalanceUseCase_Execute_CreatesMissingAccount(t *testing.T) {
	accountFrom, _ := entity.NewAccount("from", 1000)

	accountMock := &mocks.AccountGatewayMock{}
	accountMock.On("FindByID", "from").Return(accountFrom, nil)
	accountMock.On("FindByID", "new").Return(nil, entity.ErrAccountNotFound)
	accountMock.On("UpdateBalance", accountFrom).Return(nil)
	accountMock.On("Save", mock.MatchedBy(func(a *entity.Account) bool {
		return a.ID == "new" && a.Balance == 100
	})).Return(nil)

	uc := NewUpdateBalanceUseCase(accountMock)
	err := uc.Execute(UpdateBalanceInputDTO{
		AccountIDFrom:        "from",
		AccountIDTo:          "new",
		BalanceAccountIDFrom: 900,
		BalanceAccountIDTo:   100,
	})

	assert.Nil(t, err)
	accountMock.AssertExpectations(t)
}

func TestUpdateBalanceUseCase_Execute_IsIdempotent(t *testing.T) {
	accountFrom, _ := entity.NewAccount("from", 1000)
	accountTo, _ := entity.NewAccount("to", 1000)

	accountMock := &mocks.AccountGatewayMock{}
	accountMock.On("FindByID", "from").Return(accountFrom, nil)
	accountMock.On("FindByID", "to").Return(accountTo, nil)
	accountMock.On("UpdateBalance", mock.Anything).Return(nil)

	uc := NewUpdateBalanceUseCase(accountMock)
	input := UpdateBalanceInputDTO{
		AccountIDFrom:        "from",
		AccountIDTo:          "to",
		BalanceAccountIDFrom: 900,
		BalanceAccountIDTo:   1100,
	}

	assert.Nil(t, uc.Execute(input))
	assert.Nil(t, uc.Execute(input)) // mesma mensagem entregue duas vezes

	assert.Equal(t, float64(900), accountFrom.Balance)
	assert.Equal(t, float64(1100), accountTo.Balance)
}

func TestUpdateBalanceUseCase_Execute_InvalidAccountID(t *testing.T) {
	accountMock := &mocks.AccountGatewayMock{}
	accountMock.On("FindByID", "").Return(nil, entity.ErrAccountNotFound)

	uc := NewUpdateBalanceUseCase(accountMock)
	err := uc.Execute(UpdateBalanceInputDTO{
		AccountIDFrom:        "",
		AccountIDTo:          "to",
		BalanceAccountIDFrom: 900,
		BalanceAccountIDTo:   1100,
	})

	assert.ErrorIs(t, err, entity.ErrAccountIDRequired)
	accountMock.AssertNotCalled(t, "Save", mock.Anything)
}

func TestUpdateBalanceUseCase_Execute_GatewayError(t *testing.T) {
	accountMock := &mocks.AccountGatewayMock{}
	accountMock.On("FindByID", "from").Return(nil, errors.New("db down"))

	uc := NewUpdateBalanceUseCase(accountMock)
	err := uc.Execute(UpdateBalanceInputDTO{
		AccountIDFrom:        "from",
		AccountIDTo:          "to",
		BalanceAccountIDFrom: 900,
		BalanceAccountIDTo:   1100,
	})

	assert.EqualError(t, err, "db down")
	accountMock.AssertNotCalled(t, "Save", mock.Anything)
	accountMock.AssertNotCalled(t, "UpdateBalance", mock.Anything)
}

func TestUpdateBalanceInputDTO_MatchesWalletPayload(t *testing.T) {
	walletPayload := `{"account_id_from":"from","account_id_to":"to","balance_account_id_from":900,"balance_account_id_to":1100}`

	var input UpdateBalanceInputDTO
	err := json.Unmarshal([]byte(walletPayload), &input)

	assert.NoError(t, err)
	assert.Equal(t, "from", input.AccountIDFrom)
	assert.Equal(t, "to", input.AccountIDTo)
	assert.Equal(t, float64(900), input.BalanceAccountIDFrom)
	assert.Equal(t, float64(1100), input.BalanceAccountIDTo)
}
