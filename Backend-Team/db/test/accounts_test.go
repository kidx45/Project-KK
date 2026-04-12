package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	sqlc "github.com/kidx45/Project-KK/Backend-Team/db/sqlc"
)

func TestCreateAccount(t *testing.T) {
	arg := sqlc.CreateAccountParams{
		Username: "kidx45",
		Balance:  1000,
		Currency: "USD",
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Username, account.Username)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
}