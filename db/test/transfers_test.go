package db

import (
	"context"
	"testing"

	db "github.com/kidx45/Project-KK/Backend-Team/db/sqlc"
	"github.com/kidx45/Project-KK/Backend-Team/utils"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, account1 db.Account, account2 db.Account) (db.Transfer, error) {
	
	transfer, err := testQueries.CreateTransfer(context.Background(), db.CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        utils.RandomInt(10,10),
	})
	return transfer, err
}

func TestCreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	transfer, err := createRandomTransfer(t, account1, account2)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, account1.ID, transfer.FromAccountID)
	require.Equal(t, account2.ID, transfer.ToAccountID)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
}
func TestListTransfers(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	transfersBefore := []db.Transfer{}
	for i := 0; i < 5; i++ {
		transfer, err := createRandomTransfer(t, account1, account2)
		require.NoError(t, err)
		require.NotEmpty(t, transfer)
		transfersBefore = append(transfersBefore, transfer)
	}
	transfersAfter, err := testQueries.ListTransfers(context.Background(), db.ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Limit:         5,
		Offset:        0,
	})
	require.NoError(t, err)
	require.NotEmpty(t, transfersAfter)
	for i := 0; i < len(transfersBefore); i++ {
		transferBefore := transfersBefore[i]
		transferAfter := transfersAfter[i]
		require.Equal(t, transferBefore.ID, transferAfter.ID)
		require.Equal(t, transferBefore.FromAccountID, transferAfter.FromAccountID)
		require.Equal(t, transferBefore.ToAccountID, transferAfter.ToAccountID)
		require.Equal(t, transferBefore.Amount, transferAfter.Amount)
		require.NotEmpty(t, transferAfter)
		require.NotZero(t, transferAfter.ID)
		require.NotZero(t, transferAfter.CreatedAt)
		require.Equal(t, account1.ID, transferAfter.FromAccountID)
		require.Equal(t, account2.ID, transferAfter.ToAccountID)

	}
}
func TestGetTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	transferCreated, err := createRandomTransfer(t, account1, account2)
	require.NoError(t, err)
	transfer, err := testQueries.GetTransfer(context.Background(), db.GetTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, transferCreated.ID, transfer.ID)
	require.Equal(t, transferCreated.Amount, transfer.Amount)
	require.Equal(t, account1.ID, transfer.FromAccountID)
	require.Equal(t, account2.ID, transfer.ToAccountID)
}