package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// * TestTransferTx
//     * create a store with testDb
//     * create two test account
//     * concurrently exec 5 transfer tx
//       * assert no errors and results no empty
//       * after these tx, assert account1 and account2 balance changes
func TestTransferTx(t *testing.T) {
	store := NewStore(testDb)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	results := make(chan TransferTxResult)
	errs := make(chan error)
	n := 5
	amount := int64(10)	

	for i := 0; i < n; i += 1 {
		go func ()  {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID: account2.ID,
				Amount: amount,
			})
			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i += 1 {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)
		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.Equal(t, transfer.Amount, amount)
		// check fromEntry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.Equal(t, fromEntry.Amount, -amount)
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, toEntry.Amount, amount)
		// check accounts
		fromAccount := result.FromAccount
		toAccount := result.ToAccount
		require.NotEmpty(t, fromAccount)
		require.NotEmpty(t, toAccount)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.Equal(t, diff1 / amount, i + 1)
	}

	// get newest accounts and assert balance change
	latestAccount1, _ := testQueries.GetAccount(context.Background(), account1.ID)
	latestAccount2, _ := testQueries.GetAccount(context.Background(), account2.ID)
	require.Equal(t, account1.Balance - latestAccount1.Balance, amount * int64(n))
	require.Equal(t, latestAccount2.Balance - account2.Balance, amount * int64(n))
}
