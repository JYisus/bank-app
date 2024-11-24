package service_test

import (
	"context"
	"log/slog"
	"math/rand"
	"os"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/jyisus/bank-server/internal"
	"github.com/jyisus/bank-server/internal/memrepo"
	"github.com/jyisus/bank-server/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountsService_CreateAccount(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		owner          string
		initialBalance float32
		expectedError  error
	}{
		"Account successfully created": {
			owner:          faker.Name(),
			initialBalance: rand.Float32() * 100,
			expectedError:  nil,
		},
		"Invalid name": {
			owner:          "123",
			initialBalance: rand.Float32() * 100,
			expectedError:  &internal.ErrInvalidValue{},
		},
		"Invalid initial balance": {
			owner:          faker.Name(),
			initialBalance: -100,
			expectedError:  &internal.ErrInvalidValue{},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			var (
				accountsRepo    = memrepo.NewAccountsRepository()
				logger          = slog.New(slog.NewTextHandler(os.Stdout, nil))
				accountsService = service.NewAccountService(logger, accountsRepo)
				ctx             = context.Background()
			)

			actualAccount, err := accountsService.CreateAccount(ctx, tc.owner, tc.initialBalance)
			if err != nil {
				require.ErrorAs(t, err, tc.expectedError)
				return
			}

			assert.Equal(t, tc.owner, string(actualAccount.Owner))
			assert.Equal(t, tc.initialBalance, actualAccount.Balance)

			repoAccount, err := accountsRepo.Get(ctx, actualAccount.ID)
			require.NoError(t, err)

			assert.Equal(t, actualAccount, repoAccount)
		})
	}
}

func TestAccountsService_GetAccount(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		account       *internal.Account
		accountExists bool
		expectedError error
	}{
		"Account exists": {
			account:       fakeAccount(t),
			accountExists: true,
			expectedError: nil,
		},
		"Account doesn't exist": {
			account:       fakeAccount(t),
			accountExists: false,
			expectedError: &internal.ErrAccountNotFound{},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			var (
				accountsRepo    = memrepo.NewAccountsRepository()
				logger          = slog.New(slog.NewTextHandler(os.Stdout, nil))
				accountsService = service.NewAccountService(logger, accountsRepo)
				ctx             = context.Background()
			)

			if tc.accountExists {
				err := accountsRepo.Create(ctx, *tc.account)
				require.NoError(t, err)
			}

			actualAccount, err := accountsService.GetAccount(ctx, tc.account.ID)
			if err != nil {
				require.ErrorAs(t, err, tc.expectedError)
				return
			}

			assert.Equal(t, tc.account, actualAccount)
		})
	}
}

func TestAccountsService_ListAccounts(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		accounts      []internal.Account
		expectedError error
	}{
		"Multiple accounts exists": {
			accounts: []internal.Account{*fakeAccount(t), *fakeAccount(t)},
		},
		"No accounts": {
			accounts: nil,
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			var (
				accountsRepo    = memrepo.NewAccountsRepository()
				logger          = slog.New(slog.NewTextHandler(os.Stdout, nil))
				accountsService = service.NewAccountService(logger, accountsRepo)
				ctx             = context.Background()
			)

			for _, account := range tc.accounts {
				require.NoError(t, accountsRepo.Create(ctx, account))
			}

			actualAccounts, err := accountsService.ListAccounts(ctx)
			require.NoError(t, err)

			if len(tc.accounts) == 0 {
				assert.Empty(t, actualAccounts)
			} else {
				assert.ElementsMatch(t, tc.accounts, actualAccounts)
			}
		})
	}
}

func TestAccountsService_Transfer_OK(t *testing.T) {
	var (
		accountsRepo    = memrepo.NewAccountsRepository()
		logger          = slog.New(slog.NewTextHandler(os.Stdout, nil))
		accountsService = service.NewAccountService(logger, accountsRepo)
		ctx             = context.Background()

		sourceAccount = &internal.Account{
			ID:      uuid.NewString(),
			Balance: 100,
		}
		destinationAccount = &internal.Account{
			ID:      uuid.NewString(),
			Balance: 100,
		}
		transferAmount = float32(100)
	)

	require.NoError(t, accountsRepo.Create(ctx, *sourceAccount))
	require.NoError(t, accountsRepo.Create(ctx, *destinationAccount))

	err := accountsService.Transfer(
		ctx,
		sourceAccount.ID,
		destinationAccount.ID,
		transferAmount,
	)
	require.NoError(t, err)

	actualSourceAccount, err := accountsRepo.Get(ctx, sourceAccount.ID)
	require.NoError(t, err)

	actualDstAccount, err := accountsRepo.Get(ctx, destinationAccount.ID)
	require.NoError(t, err)

	expectedSourceBalance := sourceAccount.Balance - transferAmount
	expectedDstBalance := destinationAccount.Balance + transferAmount

	assert.Equal(t, expectedSourceBalance, actualSourceAccount.Balance)
	assert.Equal(t, expectedDstBalance, actualDstAccount.Balance)
}

func TestAccountsService_Transfer_SourceAccountDoesntExist(t *testing.T) {
	var (
		accountsRepo    = memrepo.NewAccountsRepository()
		logger          = slog.New(slog.NewTextHandler(os.Stdout, nil))
		accountsService = service.NewAccountService(logger, accountsRepo)
		ctx             = context.Background()

		sourceAccountID    = uuid.NewString()
		destinationAccount = &internal.Account{
			ID:      uuid.NewString(),
			Balance: 100,
		}
		transferAmount = float32(100)
	)

	require.NoError(t, accountsRepo.Create(ctx, *destinationAccount))

	err := accountsService.Transfer(
		ctx,
		sourceAccountID,
		destinationAccount.ID,
		transferAmount,
	)
	require.ErrorAs(t, err, &internal.ErrAccountNotFound{})
}

func TestAccountsService_Transfer_DstAccountDoesntExist(t *testing.T) {
	var (
		accountsRepo    = memrepo.NewAccountsRepository()
		logger          = slog.New(slog.NewTextHandler(os.Stdout, nil))
		accountsService = service.NewAccountService(logger, accountsRepo)
		ctx             = context.Background()

		sourceAccount = &internal.Account{
			ID:      uuid.NewString(),
			Balance: 100,
		}
		destinationAccountID = uuid.NewString()
		transferAmount       = float32(100)
	)

	require.NoError(t, accountsRepo.Create(ctx, *sourceAccount))

	err := accountsService.Transfer(
		ctx,
		sourceAccount.ID,
		destinationAccountID,
		transferAmount,
	)
	require.ErrorAs(t, err, &internal.ErrAccountNotFound{})
}

func TestAccountsService_Transfer_InsufficientFounts(t *testing.T) {
	var (
		accountsRepo    = memrepo.NewAccountsRepository()
		logger          = slog.New(slog.NewTextHandler(os.Stdout, nil))
		accountsService = service.NewAccountService(logger, accountsRepo)
		ctx             = context.Background()

		sourceAccount = &internal.Account{
			ID:      uuid.NewString(),
			Balance: 100,
		}
		destinationAccount = &internal.Account{
			ID:      uuid.NewString(),
			Balance: 100,
		}
		transferAmount = float32(200)
	)

	require.NoError(t, accountsRepo.Create(ctx, *sourceAccount))
	require.NoError(t, accountsRepo.Create(ctx, *destinationAccount))

	err := accountsService.Transfer(
		ctx,
		sourceAccount.ID,
		destinationAccount.ID,
		transferAmount,
	)
	require.ErrorAs(t, err, &internal.ErrInsufficientBalance{})

	actualSourceAccount, err := accountsRepo.Get(ctx, sourceAccount.ID)
	require.NoError(t, err)

	actualDstAccount, err := accountsRepo.Get(ctx, destinationAccount.ID)
	require.NoError(t, err)

	assert.Equal(t, sourceAccount.Balance, actualSourceAccount.Balance)
	assert.Equal(t, destinationAccount.Balance, actualDstAccount.Balance)
}

func fakeAccount(t *testing.T) *internal.Account {
	t.Helper()
	account, err := internal.NewAccount(uuid.NewString(), faker.Name(), rand.Float32()*100)
	require.NoError(t, err)

	return &account
}
