package internal

import "context"

type AccountsRepository interface {
	Create(ctx context.Context, account Account) error
	Get(ctx context.Context, id string) (*Account, error)
	List(ctx context.Context) ([]Account, error)
	UpdateBalance(ctx context.Context, accountID string, newBalance float32) error
}

type TransactionsRepository interface {
	Save(ctx context.Context, transaction Transaction) error
	Get(ctx context.Context, id string) (Transaction, error)
	FindAllByAccount(ctx context.Context, accountID string) ([]Transaction, error)
}
