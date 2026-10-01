package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	account, err := NewAccount("a1", 1000)
	assert.Nil(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, "a1", account.ID)
	assert.Equal(t, float64(1000), account.Balance)
	assert.False(t, account.UpdatedAt.IsZero())
}

func TestNewAccountWithEmptyID(t *testing.T) {
	account, err := NewAccount("", 1000)
	assert.Nil(t, account)
	assert.ErrorIs(t, err, ErrAccountIDRequired)
}

func TestNewAccountWithZeroBalance(t *testing.T) {
	account, err := NewAccount("a1", 0)
	assert.Nil(t, err)
	assert.Equal(t, float64(0), account.Balance)
}

func TestUpdateBalance(t *testing.T) {
	account, _ := NewAccount("a1", 1000)
	before := account.UpdatedAt
	time.Sleep(time.Millisecond)

	account.UpdateBalance(900)

	assert.Equal(t, float64(900), account.Balance)
	assert.True(t, account.UpdatedAt.After(before))
}
