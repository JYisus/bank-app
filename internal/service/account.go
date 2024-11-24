package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/jyisus/bank-server/internal"
)

type AccountService struct {
	logger             *slog.Logger
	accountsRepository internal.AccountsRepository
}

func NewAccountService(
	logger *slog.Logger,
	accountsRepository internal.AccountsRepository,
) *AccountService {
	return &AccountService{
		logger:             logger,
		accountsRepository: accountsRepository,
	}
}

func (s AccountService) CreateAccount(ctx context.Context, owner string, initialBalance float32) (*internal.Account, error) {
	// Check if initial balance is valid (> 0, only 2 decimals)
	id := uuid.NewString()

	if err := s.checkIfAccountExists(ctx, id); err != nil {
		return nil, err
	}

	account, err := internal.NewAccount(id, owner, initialBalance)
	if err != nil {
		return nil, err
	}

	if err := s.accountsRepository.Create(ctx, account); err != nil {
		return nil, err
	}

	s.logger.Debug("New account created", "ID", id, "owner", owner, "initial_balance", initialBalance)

	return &account, nil
}

func (s AccountService) GetAccount(ctx context.Context, id string) (*internal.Account, error) {
	// Check if initial balance is valid (> 0, only 2 decimals)
	account, err := s.accountsRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	s.logger.Debug("Returning account", "ID", id, "owner", account.Owner, "balance", account.Balance)

	return account, nil
}

func (s AccountService) ListAccounts(ctx context.Context) ([]internal.Account, error) {
	accounts, err := s.accountsRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Debug("Returning accounts list", "totalAccounts", len(accounts))

	return accounts, nil
}

func (s AccountService) Transfer(ctx context.Context, sourceAccountID, destinationAccountID string, amount float32) error {
	sourceAccount, err := s.accountsRepository.Get(ctx, sourceAccountID)
	if err != nil {
		return fmt.Errorf("getting source account: %w", err)
	}

	destinationAccount, err := s.accountsRepository.Get(ctx, destinationAccountID)
	if err != nil {
		return fmt.Errorf("getting destination account: %w", err)
	}

	if err := sourceAccount.Withdraw(amount); err != nil {
		return err
	}

	destinationAccount.Deposit(amount)

	wg := sync.WaitGroup{}
	wg.Add(2)

	var errs error

	go func() {
		defer wg.Done()
		err = s.accountsRepository.UpdateBalance(ctx, sourceAccount.ID, sourceAccount.Balance)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("updating source account balance: %w", err))
		}
	}()

	go func() {
		defer wg.Done()
		err = s.accountsRepository.UpdateBalance(ctx, destinationAccount.ID, destinationAccount.Balance)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("updating destination account balance: %w", err))
		}
	}()

	wg.Wait()

	return errs
}

func (s AccountService) checkIfAccountExists(ctx context.Context, id string) error {
	_, err := s.accountsRepository.Get(ctx, id)
	switch {
	case errors.As(err, &internal.ErrAccountNotFound{}):
		return nil
	case err != nil:
		return err
	default:
		return internal.ErrAccountAlreadyExists
	}
}
