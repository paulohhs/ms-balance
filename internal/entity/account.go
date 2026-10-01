package entity

import (
	"errors"
	"time"
)

var (
	ErrAccountIDRequired = errors.New("account id is required")
	ErrAccountNotFound   = errors.New("account not found")
)

type Account struct {
	ID        string
	Balance   float64
	UpdatedAt time.Time
}

func NewAccount(id string, balance float64) (*Account, error) {
	account := &Account{
		ID:        id,
		Balance:   balance,
		UpdatedAt: time.Now(),
	}

	err := account.Validate()
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (a *Account) Validate() error {
	if a.ID == "" {
		return ErrAccountIDRequired
	}

	return nil
}

func (a *Account) UpdateBalance(balance float64) {
	a.Balance = balance
	a.UpdatedAt = time.Now()
}
