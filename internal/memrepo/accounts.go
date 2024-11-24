package memrepo

import (
	"context"
	"errors"
	"sync"

	"github.com/jyisus/bank-server/internal"
)

type AccountsRepository struct {
	memAccounts map[string]internal.Account
	mutex       *sync.Mutex
}

var _ internal.AccountsRepository = (*AccountsRepository)(nil)

func NewAccountsRepository() *AccountsRepository {
	return &AccountsRepository{
		memAccounts: make(map[string]internal.Account),
		mutex:       &sync.Mutex{},
	}
}

func (ar *AccountsRepository) Create(_ context.Context, account internal.Account) error {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()

	if _, ok := ar.memAccounts[account.ID]; ok {
		return internal.ErrAccountAlreadyExists
	}
	ar.memAccounts[account.ID] = account

	return nil
}

func (ar *AccountsRepository) Get(_ context.Context, id string) (*internal.Account, error) {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()

	account, ok := ar.memAccounts[id]
	if !ok {
		return nil, internal.ErrAccountNotFound{AccountID: id}
	}

	return &account, nil
}

func (ar *AccountsRepository) List(_ context.Context) ([]internal.Account, error) {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()

	accounts := make([]internal.Account, 0, len(ar.memAccounts))
	for _, account := range ar.memAccounts {
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (ar *AccountsRepository) UpdateBalance(_ context.Context, accountID string, newBalance float32) error {
	account, ok := ar.memAccounts[accountID]
	if !ok {
		return errors.New("account with given ID not found")
	}

	ar.mutex.Lock()
	defer ar.mutex.Unlock()

	account.Balance = newBalance
	ar.memAccounts[accountID] = account

	return nil
}
