package service_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/google/uuid"
	"github.com/jyisus/bank-server/internal"
	"github.com/jyisus/bank-server/internal/memrepo"
	"github.com/jyisus/bank-server/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionsService_SaveTransaction(t *testing.T) {
	t.Parallel()
	monkey.Patch(time.Now, func() time.Time { return time.Time{} })
	defer monkey.UnpatchAll()

	testCases := map[string]struct {
		account       *internal.Account
		txType        string
		amount        float32
		expectedError error
	}{
		"Deposit transaction successfully created": {
			account:       fakeAccount(t),
			txType:        internal.TxDeposit,
			amount:        100,
			expectedError: nil,
		},
		"Withdrawal transaction successfully created": {
			account: &internal.Account{
				ID:      uuid.NewString(),
				Balance: 100,
			},
			txType:        internal.TxWithdrawal,
			amount:        50,
			expectedError: nil,
		},
		"Withdrawal transaction insufficient balance": {
			account: &internal.Account{
				ID:      uuid.NewString(),
				Balance: 100,
			},
			txType:        internal.TxWithdrawal,
			amount:        2000,
			expectedError: &internal.ErrInsufficientBalance{},
		},
		"Account doesn't exist": {
			account:       nil,
			txType:        internal.TxDeposit,
			expectedError: &internal.ErrAccountNotFound{},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			var (
				logger              = slog.New(slog.NewTextHandler(os.Stdout, nil))
				accountsRepo        = memrepo.NewAccountsRepository()
				transactionsRepo    = memrepo.NewTransactionsRepository()
				transactionsService = service.NewTransactionService(logger, accountsRepo, transactionsRepo)
				ctx                 = context.Background()
			)

			accountID := ""
			if tc.account != nil {
				accountID = tc.account.ID
				require.NoError(t, accountsRepo.Create(ctx, *tc.account))
			}

			createdTx, err := transactionsService.SaveTransaction(ctx, accountID, tc.txType, tc.amount)
			if err != nil {
				assert.ErrorAs(t, err, tc.expectedError)
				return
			}

			require.NoError(t, err)

			repoTransaction, err := transactionsRepo.Get(ctx, createdTx.ID)
			require.NoError(t, err)

			assert.Equal(t, *createdTx, repoTransaction)

			account, err := accountsRepo.Get(ctx, tc.account.ID)
			require.NoError(t, err)

			if tc.txType == internal.TxDeposit {
				assert.Equal(t, tc.account.Balance+tc.amount, account.Balance)
				return
			}

			if tc.txType == internal.TxWithdrawal {
				assert.Equal(t, tc.account.Balance-tc.amount, account.Balance)
			}
		})
	}
}

func TestTransactionsService_RetrieveAccountTransactions(t *testing.T) {
	t.Parallel()
	monkey.Patch(time.Now, func() time.Time { return time.Time{} })
	defer monkey.UnpatchAll()

	testCases := map[string]struct {
		account              *internal.Account
		allTransactions      []internal.Transaction
		expectedTransactions []internal.Transaction
		expectedError        error
	}{
		"Deposit transaction successfully created": {
			account: &internal.Account{
				ID: "test-account",
			},
			allTransactions: []internal.Transaction{
				{
					ID:        "another-test-transaction",
					AccountID: "another-test-account",
				},
				{
					ID:        "test-transaction",
					AccountID: "test-account",
				},
			},
			expectedTransactions: []internal.Transaction{
				{
					ID:        "test-transaction",
					AccountID: "test-account",
				},
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			var (
				logger              = slog.New(slog.NewTextHandler(os.Stdout, nil))
				accountsRepo        = memrepo.NewAccountsRepository()
				transactionsRepo    = memrepo.NewTransactionsRepository()
				transactionsService = service.NewTransactionService(logger, accountsRepo, transactionsRepo)
				ctx                 = context.Background()
			)

			accountID := ""
			if tc.account != nil {
				accountID = tc.account.ID
				require.NoError(t, accountsRepo.Create(ctx, *tc.account))
			}

			for _, tx := range tc.allTransactions {
				require.NoError(t, transactionsRepo.Save(ctx, tx))
			}

			transactions, err := transactionsService.RetrieveAccountTransactions(ctx, accountID)
			if err != nil {
				assert.ErrorAs(t, err, tc.expectedError)
				return
			}
			require.NoError(t, err)

			assert.ElementsMatch(t, tc.expectedTransactions, transactions)
		})
	}
}
