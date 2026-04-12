package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	sqlc "github.com/kidx45/Project-KK/Backend-Team/db/sqlc"
	utils "github.com/kidx45/Project-KK/Backend-Team/utils"
)

func createRandomAccount(t *testing.T) sqlc.Account {
	user := createRandomUser(t)
	arg := sqlc.CreateAccountParams{
		Username: user.Username,
		Balance:  utils.RandomInt(1000, 1000000),
		Currency: utils.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Username, account.Username)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}
