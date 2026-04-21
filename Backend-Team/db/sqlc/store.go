package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type Store struct {
	*Queries //composition
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) execTx (ctx context.Context, fn func (*Queries) error) error {
	tx, err := s.db.BeginTx(ctx,nil)
	if err != nil {
		return err
	}	

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rollBackerror := tx.Rollback(); rollBackerror != nil {
			return fmt.Errorf("transaction err: %s, query err: %s",rollBackerror,err)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64     `json:"fromAccountId"`
	ToAccountID   int64     `json:"toAccountId"`
	Amount        int64     `json:"amount"`
}

type TransferTxResult struct {
	FromAccount Account `json:"fromAccount"`
	ToAccount Account `json:"toAccount"`
	FromEntry Entry `json:"FromEntry"`
	ToEntry Entry `json:"ToEntry"`
	Transfer Transfer `json:"transfer"`
}

var txKey = struct {}{}

func (s *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {

		var err error
		// "Hey Context, give me the Value attached to the txKey label"
		myName := ctx.Value(txKey) 

		log.Printf("%s is doing a transfer",myName)
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		account1, err := q.GetAccountByIdForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}

		result.FromAccount, err = q.UpdateAccountById(ctx, UpdateAccountByIdParams{
			ID: account1.ID,
			Balance: account1.Balance - arg.Amount,
		})
		if err != nil {
			return err
		}

		account2, err := q.GetAccountByIdForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}

		result.ToAccount, err = q.UpdateAccountById(ctx, UpdateAccountByIdParams{
			ID: account2.ID,
			Balance: account2.Balance + arg.Amount,
		})
		if err != nil {
			return err
		}

		return err
	})
	return result, err
}