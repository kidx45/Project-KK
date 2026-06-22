package db

import (
	"context"
	"fmt"
	"log"
	"testing"

	sqlc "github.com/kidx45/Project-KK/Backend-Team/db/sqlc"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := sqlc.NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	var txKey = struct{}{}

	log.Printf("Before: account1 balance: %d, account2 balance: %d", account1.Balance, account2.Balance)
	n := 3
	amount := int64(10)

	results := make(chan sqlc.TransferTxResult)
	errors := make(chan error)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, sqlc.TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errors <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)

		result := <-results
		log.Printf("Tx %d: from account balance: %d, to account balance: %d", i+1, result.FromAccount.Balance, result.ToAccount.Balance)
		require.NotEmpty(t, result)
		require.Equal(t, account1.ID, result.Transfer.FromAccountID)
		require.Equal(t, account2.ID, result.Transfer.ToAccountID)
		require.Equal(t, amount, result.Transfer.Amount)

		require.NotEmpty(t, result.FromEntry)
		require.Equal(t, result.FromEntry.AccountID, account1.ID)
		require.Equal(t, result.ToEntry.AccountID, account2.ID)
		require.Equal(t, result.FromEntry.Amount, -amount)
		require.Equal(t, result.ToEntry.Amount, amount)

		diff1 := account1.Balance - result.FromAccount.Balance
		diff2 := result.ToAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAccount1, err := testQueries.GetAccountByUsername(context.Background(), account1.Username)
	require.NoError(t, err)
	updatedAccount2, err := testQueries.GetAccountByUsername(context.Background(), account2.Username)
	require.NoError(t, err)
	log.Printf("After: account1 balance: %d, account2 balance: %d", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

func TestTransfersTx(t *testing.T) {
	store := sqlc.NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	var txKey = struct{}{}

	log.Printf("Before: account1 balance: %d, account2 balance: %d", account1.Balance, account2.Balance)
	n := 10
	amount := int64(10)

	errors := make(chan error)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i)
		FromAccountId := account1.ID
		ToAccountId := account2.ID

		if i%2 == 1 {
			FromAccountId = account2.ID
			ToAccountId = account1.ID
		}
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			_, err := store.TransferTx(ctx, sqlc.TransferTxParams{
				FromAccountID: FromAccountId,
				ToAccountID:   ToAccountId,
				Amount:        amount,
			})
			errors <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)
	}

	updatedAccount1, err := testQueries.GetAccountByUsername(context.Background(), account1.Username)
	require.NoError(t, err)
	updatedAccount2, err := testQueries.GetAccountByUsername(context.Background(), account2.Username)
	require.NoError(t, err)
	log.Printf("After: account1 balance: %d, account2 balance: %d", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
