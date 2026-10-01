package get_balance

import (
	"time"

	"github.com/paulohhs/ms-balance/internal/gateway"
)

type GetBalanceInputDTO struct {
	AccountID string `json:"account_id"`
}

type GetBalanceOutputDTO struct {
	AccountID string    `json:"account_id"`
	Balance   float64   `json:"balance"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetBalanceUseCase struct {
	AccountGateway gateway.AccountGateway
}

func NewGetBalanceUseCase(accountGateway gateway.AccountGateway) *GetBalanceUseCase {
	return &GetBalanceUseCase{
		AccountGateway: accountGateway,
	}
}

func (uc *GetBalanceUseCase) Execute(input GetBalanceInputDTO) (*GetBalanceOutputDTO, error) {
	account, err := uc.AccountGateway.FindByID(input.AccountID)
	if err != nil {
		return nil, err
	}

	return &GetBalanceOutputDTO{
		AccountID: account.ID,
		Balance:   account.Balance,
		UpdatedAt: account.UpdatedAt,
	}, nil
}
