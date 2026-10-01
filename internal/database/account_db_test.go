package database

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/paulohhs/ms-balance/internal/entity"
	"github.com/stretchr/testify/suite"
)

type AccountDBTestSuite struct {
	suite.Suite
	db        *sql.DB
	accountDB *AccountDB
}

func (s *AccountDBTestSuite) SetupTest() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Require().NoError(err)
	db.SetMaxOpenConns(1)

	_, err = db.Exec("CREATE TABLE accounts (id varchar(255) PRIMARY KEY, balance float, updated_at datetime)")
	s.Require().NoError(err)

	s.db = db
	s.accountDB = NewAccountDB(db)
}

func (s *AccountDBTestSuite) TearDownTest() {
	s.db.Close()
}

func TestAccountDBTestSuite(t *testing.T) {
	suite.Run(t, new(AccountDBTestSuite))
}

func (s *AccountDBTestSuite) TestSchemaIsCreated() {
	var name string
	err := s.db.QueryRow("SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'accounts'").Scan(&name)
	s.NoError(err)
	s.Equal("accounts", name)
}

func (s *AccountDBTestSuite) TestSave() {
	account, _ := entity.NewAccount("a1", 1000)

	err := s.accountDB.Save(account)
	s.NoError(err)

	var balance float64
	err = s.db.QueryRow("SELECT balance FROM accounts WHERE id = ?", "a1").Scan(&balance)
	s.NoError(err)
	s.Equal(float64(1000), balance)
}

func (s *AccountDBTestSuite) TestSaveWithDuplicatedID() {
	account, _ := entity.NewAccount("a1", 1000)
	s.NoError(s.accountDB.Save(account))

	err := s.accountDB.Save(account)
	s.Error(err)
}

func (s *AccountDBTestSuite) TestFindByID() {
	account, _ := entity.NewAccount("a1", 1000)
	s.NoError(s.accountDB.Save(account))

	accountDB, err := s.accountDB.FindByID("a1")
	s.NoError(err)
	s.Equal(account.ID, accountDB.ID)
	s.Equal(account.Balance, accountDB.Balance)
	s.WithinDuration(account.UpdatedAt, accountDB.UpdatedAt, time.Second)
}

func (s *AccountDBTestSuite) TestFindByIDNotFound() {
	accountDB, err := s.accountDB.FindByID("does-not-exist")
	s.Nil(accountDB)
	s.ErrorIs(err, entity.ErrAccountNotFound)
}

func (s *AccountDBTestSuite) TestUpdateBalance() {
	account, _ := entity.NewAccount("a1", 1000)
	s.NoError(s.accountDB.Save(account))

	account.UpdateBalance(750)
	err := s.accountDB.UpdateBalance(account)
	s.NoError(err)

	accountDB, err := s.accountDB.FindByID("a1")
	s.NoError(err)
	s.Equal(float64(750), accountDB.Balance)
}
