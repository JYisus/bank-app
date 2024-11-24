package internal_test

import (
	"math/rand"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/jyisus/bank-server/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccount_New_OK(t *testing.T) {
	var (
		id     = uuid.NewString()
		owner  = faker.Name()
		amount = rand.Float32() * 100
	)

	_, err := internal.NewAccount(id, owner, amount)
	require.NoError(t, err)
}

func TestAccount_New_InvalidOwner(t *testing.T) {
	var (
		id     = uuid.NewString()
		owner  = "12345"
		amount = rand.Float32()
	)

	_, err := internal.NewAccount(id, owner, amount)
	require.ErrorAs(t, err, &internal.ErrInvalidValue{})
}

func TestAccount_New_InvalidAmount(t *testing.T) {
	var (
		id     = uuid.NewString()
		owner  = faker.Name()
		amount = -rand.Float32()
	)

	_, err := internal.NewAccount(id, owner, amount)
	require.ErrorAs(t, err, &internal.ErrInvalidValue{})
}

func TestAccount_Deposit(t *testing.T) {
	var (
		id     = uuid.NewString()
		owner  = "Test User"
		amount = 100
	)

	account, err := internal.NewAccount(id, owner, float32(amount))
	require.NoError(t, err)

	account.Deposit(100)
	assert.Equal(t, float32(200), account.Balance)
}

func TestAccount_Withdrawal_OK(t *testing.T) {
	var (
		id     = uuid.NewString()
		owner  = "Test User"
		amount = 100
	)

	account, err := internal.NewAccount(id, owner, float32(amount))
	require.NoError(t, err)

	err = account.Withdraw(20)

	require.NoError(t, err)
	assert.Equal(t, float32(80), account.Balance)
}

func TestAccount_Withdrawal_InsufficentBalance(t *testing.T) {
	var (
		id     = uuid.NewString()
		owner  = "Test User"
		amount = 100
	)

	account, err := internal.NewAccount(id, owner, float32(amount))
	require.NoError(t, err)

	err = account.Withdraw(200)
	assert.ErrorAs(t, err, &internal.ErrInsufficientBalance{})
	assert.Equal(t, float32(100), account.Balance)
}
