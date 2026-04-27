package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type Store struct {
	*Queries //composition
	db       *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rollBackerror := tx.Rollback(); rollBackerror != nil {
			return fmt.Errorf("transaction err: %s, query err: %s", rollBackerror, err)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"fromAccountId"`
	ToAccountID   int64 `json:"toAccountId"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	FromAccount Account  `json:"fromAccount"`
	ToAccount   Account  `json:"toAccount"`
	FromEntry   Entry    `json:"FromEntry"`
	ToEntry     Entry    `json:"ToEntry"`
	Transfer    Transfer `json:"transfer"`
}

var txKey = struct{}{}

func (s *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {

		var err error
		// "Hey Context, give me the Value attached to the txKey label"
		myName := ctx.Value(txKey)

		log.Printf("%s is doing a transfer", myName)
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		log.Printf("%s is creating an entry for account %d with %d",myName, arg.FromAccountID,-arg.Amount)
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		log.Printf("%s is creating an entry for account %d with %d",myName, arg.ToAccountID,arg.Amount)
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			log.Printf("%s is subtracting from account %d with %d",myName, arg.FromAccountID,-arg.Amount)
			result.FromAccount, err = q.AddMoneyIntoAccount(ctx, AddMoneyIntoAccountParams{
				ID:     arg.FromAccountID,
				Amount: -arg.Amount,
			})
			if err != nil {
				return err
			}

			log.Printf("%s is adding to account %d with %d",myName, arg.ToAccountID,arg.Amount)
			result.ToAccount, err = q.AddMoneyIntoAccount(ctx, AddMoneyIntoAccountParams{
				ID:     arg.ToAccountID,
				Amount: +arg.Amount,
			})
			if err != nil {
				return err
			}
		} else {
			log.Printf("%s is adding to account %d with %d",myName, arg.ToAccountID,arg.Amount)
			result.ToAccount, err = q.AddMoneyIntoAccount(ctx, AddMoneyIntoAccountParams{
				ID:     arg.ToAccountID,
				Amount: +arg.Amount,
			})
			if err != nil {
				return err
			}

			log.Printf("%s is subtracting from account %d with %d",myName, arg.FromAccountID,-arg.Amount)
			result.FromAccount, err = q.AddMoneyIntoAccount(ctx, AddMoneyIntoAccountParams{
				ID:     arg.FromAccountID,
				Amount: -arg.Amount,
			})
			if err != nil {
				return err
			}
		}
		return err
	})
	return result, err
}
