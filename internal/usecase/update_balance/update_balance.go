package updatebalance

import (
	"errors"

	"github.com/paulohhs/ms-balance/internal/entity"
	"github.com/paulohhs/ms-balance/internal/gateway"
)

type UpdateBalanceInputDTO struct {
	AccountIDFrom        string  `json:"account_id_from"`
	AccountIDTo          string  `json:"account_id_to"`
	BalanceAccountIDFrom float64 `json:"balance_account_id_from"`
	BalanceAccountIDTo   float64 `json:"balance_account_id_to"`
}

type UpdateBalanceUseCase struct {
	AccountGateway gateway.AccountGateway
}

func NewUpdateBalanceUseCase(accountGateway gateway.AccountGateway) *UpdateBalanceUseCase {
	return &UpdateBalanceUseCase{
		AccountGateway: accountGateway,
	}
}

func (uc *UpdateBalanceUseCase) Execute(input UpdateBalanceInputDTO) error {
	err := saveBalance(uc.AccountGateway, input.AccountIDFrom, input.BalanceAccountIDFrom)
	if err != nil {
		return err
	}

	return saveBalance(uc.AccountGateway, input.AccountIDTo, input.BalanceAccountIDTo)
}

func saveBalance(accountRepository gateway.AccountGateway, accountID string, balance float64) error {
	account, err := accountRepository.FindByID(accountID)
	if errors.Is(err, entity.ErrAccountNotFound) {
		account, err = entity.NewAccount(accountID, balance)
		if err != nil {
			return err
		}
		return accountRepository.Save(account)
	}
	if err != nil {
		return err
	}

	account.UpdateBalance(balance)
	return accountRepository.UpdateBalance(account)
}
