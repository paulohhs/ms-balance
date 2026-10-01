package database

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type AccountDBTestSuite struct {
	suite.Suite
	db *sql.DB
}

func (s *AccountDBTestSuite) SetupTest() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Require().NoError(err)
	db.SetMaxOpenConns(1)

	_, err = db.Exec("CREATE TABLE accounts (id varchar(255) PRIMARY KEY, balance float, updated_at datetime)")
	s.Require().NoError(err)

	s.db = db
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
