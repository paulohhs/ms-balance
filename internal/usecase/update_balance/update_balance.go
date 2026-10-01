package update_balance

import (
	"context"
	"errors"

	"github.com/paulohhs/ms-balance/internal/entity"
	"github.com/paulohhs/ms-balance/internal/gateway"
	"github.com/paulohhs/ms-balance/pkg/uow"
)

type UpdateBalanceInputDTO struct {
	AccountIDFrom        string  `json:"account_id_from"`
	AccountIDTo          string  `json:"account_id_to"`
	BalanceAccountIDFrom float64 `json:"balance_account_id_from"`
	BalanceAccountIDTo   float64 `json:"balance_account_id_to"`
}

type UpdateBalanceUseCase struct {
	Uow uow.UowInterface
}

func NewUpdateBalanceUseCase(uow uow.UowInterface) *UpdateBalanceUseCase {
	return &UpdateBalanceUseCase{
		Uow: uow,
	}
}

func (uc *UpdateBalanceUseCase) Execute(ctx context.Context, input UpdateBalanceInputDTO) error {
	return uc.Uow.Do(ctx, func(_ *uow.Uow) error {
		accountRepository, err := uc.getAccountRepository(ctx)
		if err != nil {
			return err
		}

		err = saveBalance(accountRepository, input.AccountIDFrom, input.BalanceAccountIDFrom)
		if err != nil {
			return err
		}

		return saveBalance(accountRepository, input.AccountIDTo, input.BalanceAccountIDTo)
	})
}

func (uc *UpdateBalanceUseCase) getAccountRepository(ctx context.Context) (gateway.AccountGateway, error) {
	repo, err := uc.Uow.GetRepository(ctx, "AccountDB")
	if err != nil {
		return nil, err
	}

	accountRepository, ok := repo.(gateway.AccountGateway)
	if !ok {
		return nil, errors.New("AccountDB repository does not implement gateway.AccountGateway")
	}

	return accountRepository, nil
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
