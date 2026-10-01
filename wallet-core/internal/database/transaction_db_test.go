package database

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/paulohhs/ms-wallet/internal/entity"
	"github.com/stretchr/testify/suite"
)

type TransactionDBTestSuite struct {
	suite.Suite
	db            *sql.DB
	client        *entity.Client
	client2       *entity.Client
	accountFrom   *entity.Account
	accountTo     *entity.Account
	transactionDB *TransactionDB
}

func (t *TransactionDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	t.Nil(err)
	t.db = db
	db.Exec("Create table clients (id varchar(255), name varchar(255), email varchar(255), created_at date)")
	db.Exec("Create table accounts (id varchar(255), client_id varchar(255), balance float, created_at date)")
	db.Exec("Create table transactions (id varchar(255), account_id_from varchar(255), account_id_to varchar(255), amount float, created_at date)")

	client, _ := entity.NewClient("Test Client", "t@t.com")
	t.Nil(err)
	t.client = client
	client2, _ := entity.NewClient("Test Client", "t@t.com")
	t.Nil(err)
	t.client2 = client2

	accountFrom := entity.NewAccount(t.client)
	accountFrom.Balance = 1000
	t.accountFrom = accountFrom
	accountTo := entity.NewAccount(t.client2)
	accountTo.Balance = 1000
	t.accountTo = accountTo

	t.transactionDB = NewTransactionDB(db)
}

func (s *TransactionDBTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE transactions")
	s.db.Exec("DROP TABLE accounts")
	s.db.Exec("DROP TABLE clients")
}

func TestTransactionDBTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionDBTestSuite))
}

func (t *TransactionDBTestSuite) TestCreate() {
	transaction, err := entity.NewTransaction(t.accountFrom, t.accountTo, 100)
	t.Nil(err)

	err = t.transactionDB.Create(transaction)
	t.Nil(err)
}
