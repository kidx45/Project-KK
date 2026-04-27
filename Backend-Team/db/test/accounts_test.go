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

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)
	accountGot,err := testQueries.GetAccountByUsername(context.Background(),account.Username)
	_,err2 := testQueries.GetAccountById(context.Background(),account.ID + 1)
	require.Equal(t,accountGot.Username,account.Username)
	require.Equal(t,account.Balance,accountGot.Balance)
	require.NoError(t,err)
	require.Error(t,err2)
}

func TestListAccounts(t *testing.T) {
	accounts := make([]sqlc.Account,0)
	for i := 0; i < 5; i++ {
		accounts = append(accounts,createRandomAccount(t))
	}
	
	accountGot,err := testQueries.ListAccountsDesc(context.Background(),sqlc.ListAccountsDescParams{
		Limit: 5,
		Offset: 0,
	})
	require.NoError(t,err)
	for i := 4; i >= 0; i-- {
		require.NotEmpty(t, accountGot[4-i])
		require.Equal(t,accountGot[4-i].Username,accounts[i].Username)
		require.Equal(t,accountGot[4-i].Balance,accounts[i].Balance)
		require.Equal(t,accountGot[4-i].Currency,accounts[i].Currency)
	}
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)
	param := sqlc.UpdateAccountByIdParams{
		ID: account.ID,
		Balance: account.Balance - 100,
	}
	updatedAccount,err := testQueries.UpdateAccountById(context.Background(), param)
	require.NoError(t,err)
	require.Equal(t,updatedAccount.Balance,account.Balance-100)
}