package update_balance

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/paulohhs/ms-balance/internal/database"
	"github.com/paulohhs/ms-balance/internal/entity"
	"github.com/paulohhs/ms-balance/pkg/uow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRealUow(t *testing.T) (*sql.DB, *uow.Uow) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	_, err = db.Exec("CREATE TABLE accounts (id varchar(255) PRIMARY KEY, balance float, updated_at datetime)")
	require.NoError(t, err)

	unitOfWork := uow.NewUow(context.Background(), db)
	unitOfWork.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(tx)
	})
	return db, unitOfWork
}

func TestUpdateBalanceUseCase_Execute_PersistsBothAccounts(t *testing.T) {
	db, unitOfWork := setupRealUow(t)
	defer db.Close()

	uc := NewUpdateBalanceUseCase(unitOfWork)
	err := uc.Execute(context.Background(), UpdateBalanceInputDTO{
		AccountIDFrom:        "from",
		AccountIDTo:          "to",
		BalanceAccountIDFrom: 900,
		BalanceAccountIDTo:   1100,
	})
	require.NoError(t, err)

	var count int
	require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM accounts").Scan(&count))
	assert.Equal(t, 2, count)
}

func TestUpdateBalanceUseCase_Execute_RollsBackWhenSecondAccountFails(t *testing.T) {
	db, unitOfWork := setupRealUow(t)
	defer db.Close()

	uc := NewUpdateBalanceUseCase(unitOfWork)
	err := uc.Execute(context.Background(), UpdateBalanceInputDTO{
		AccountIDFrom:        "from",
		AccountIDTo:          "", // inválido: falha depois de gravar "from"
		BalanceAccountIDFrom: 900,
		BalanceAccountIDTo:   1100,
	})
	assert.ErrorIs(t, err, entity.ErrAccountIDRequired)

	var count int
	require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM accounts").Scan(&count))
	assert.Equal(t, 0, count, "a conta 'from' não pode ficar gravada após o rollback")
}
