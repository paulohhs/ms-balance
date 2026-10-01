package database

import (
	"database/sql"
	"errors"

	"github.com/paulohhs/ms-balance/internal/entity"
)

type AccountDB struct {
	DB *sql.DB
}

func NewAccountDB(db *sql.DB) *AccountDB {
	return &AccountDB{DB: db}
}

func (a *AccountDB) FindByID(id string) (*entity.Account, error) {
	stmt, err := a.DB.Prepare("SELECT id, balance, updated_at FROM accounts WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var account entity.Account
	err = stmt.QueryRow(id).Scan(&account.ID, &account.Balance, &account.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, entity.ErrAccountNotFound
	}
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (a *AccountDB) Save(account *entity.Account) error {
	stmt, err := a.DB.Prepare("INSERT INTO accounts (id, balance, updated_at) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(account.ID, account.Balance, account.UpdatedAt)
	return err
}

func (a *AccountDB) UpdateBalance(account *entity.Account) error {
	stmt, err := a.DB.Prepare("UPDATE accounts SET balance = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(account.Balance, account.UpdatedAt, account.ID)
	return err
}
