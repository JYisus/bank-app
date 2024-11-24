package memrepo

import (
	"context"
	"errors"

	"github.com/jyisus/bank-server/internal"
)

type TransactionsReposiory struct {
	transactions map[string]internal.Transaction
}

var _ internal.TransactionsRepository = (*TransactionsReposiory)(nil)

func NewTransactionsRepository() *TransactionsReposiory {
	return &TransactionsReposiory{
		transactions: make(map[string]internal.Transaction),
	}
}

func (tr *TransactionsReposiory) Save(_ context.Context, transaction internal.Transaction) error {
	if _, ok := tr.transactions[transaction.ID]; ok {
		return errors.New("transaction with given ID already exists")
	}

	tr.transactions[transaction.ID] = transaction

	return nil
}

func (tr *TransactionsReposiory) Get(_ context.Context, id string) (internal.Transaction, error) {
	transaction, ok := tr.transactions[id]
	if !ok {
		return internal.Transaction{}, errors.New("transaction not found")
	}

	return transaction, nil
}

func (tr *TransactionsReposiory) FindAllByAccount(_ context.Context, accountID string) ([]internal.Transaction, error) {
	transactions := make([]internal.Transaction, 0, len(tr.transactions))
	for _, transaction := range tr.transactions {
		if transaction.AccountID == accountID {
			transactions = append(transactions, transaction)
		}
	}

	return transactions, nil
}
