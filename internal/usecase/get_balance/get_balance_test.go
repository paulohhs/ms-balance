package get_balance

import (
	"testing"

	"github.com/paulohhs/ms-balance/internal/entity"
	"github.com/paulohhs/ms-balance/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetBalanceUseCase_Execute(t *testing.T) {
	account, _ := entity.NewAccount("a1", 1000)
	accountMock := &mocks.AccountGatewayMock{}
	accountMock.On("FindByID", "a1").Return(account, nil)

	uc := NewGetBalanceUseCase(accountMock)
	output, err := uc.Execute(GetBalanceInputDTO{AccountID: "a1"})

	assert.Nil(t, err)
	assert.Equal(t, "a1", output.AccountID)
	assert.Equal(t, float64(1000), output.Balance)
	assert.Equal(t, account.UpdatedAt, output.UpdatedAt)
	accountMock.AssertExpectations(t)
	accountMock.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGetBalanceUseCase_Execute_AccountNotFound(t *testing.T) {
	accountMock := &mocks.AccountGatewayMock{}
	accountMock.On("FindByID", "x").Return(nil, entity.ErrAccountNotFound)

	uc := NewGetBalanceUseCase(accountMock)
	output, err := uc.Execute(GetBalanceInputDTO{AccountID: "x"})

	assert.Nil(t, output)
	assert.ErrorIs(t, err, entity.ErrAccountNotFound)
}
