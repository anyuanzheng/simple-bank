package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	db *sql.DB
	// 封装所有Queries的功能，同时持有db instance，以实现transaction
	*Queries
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
		Queries: New(db),
	}
}

var txKey = struct{}{}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// ？ 为什么要创建一个新的Queries
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 `json:"to_account_id"`
	Amount int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`	
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var transferResult TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		transferResult.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		transferResult.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		transferResult.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		// update from account
		transferResult.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{ ID: arg.FromAccountID, Amount: -arg.Amount})
		if err != nil {
			return err
		}
		
		// update to account
		transferResult.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{ ID: arg.ToAccountID, Amount: arg.Amount})
		if err != nil {
			return err
		}	

		return nil
	})

	return transferResult, err
}