package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/iamzay/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	params := CreateAccountParams{
		Owner: user.Username,
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, params.Owner, account.Owner)
	require.Equal(t, params.Balance, account.Balance)
	require.Equal(t, account.Currency, account.Currency)
	require.NotZero(t, account.CreatedAt)
	require.NotZero(t, account.ID)
	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)

  account1, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.Equal(t, account.ID, account1.ID)
	require.Equal(t, account.Balance, account1.Balance)
	require.Equal(t, account.Currency, account1.Currency)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	account1, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Empty(t, account1)
	require.Error(t, err)
	require.EqualError(t, sql.ErrNoRows, err.Error())
}

func TestListAccounts(t *testing.T) {
	n := 5
	var lastAccount Account
	for i := 0; i < n; i += 1 {
		lastAccount = createRandomAccount(t)
	}

	accounts, err := testQueries.ListAccounts(context.Background(), ListAccountsParams{ Offset: 0, Limit: 5, Owner: lastAccount.Owner })
	require.NoError(t, err)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func TestUpdateAccount(t * testing.T) {
	account := createRandomAccount(t)

	randomMoney := util.RandomMoney()
	account1, err := testQueries.UpdateAccount(context.Background(), UpdateAccountParams{ID: account.ID, Balance: randomMoney})

	require.NoError(t, err)
	require.Equal(t, account.Owner, account1.Owner)
	require.Equal(t, account.ID, account1.ID)
	require.Equal(t, account.Currency, account1.Currency)
	require.Equal(t, account1.Balance, randomMoney)
}
