package internal

import (
	"time"
)

// TODO: use int?
type TransactionType string

const (
	TxDeposit    = "deposit"
	TxWithdrawal = "withdrawal"
)

func NewTransactionType(txType string) (TransactionType, error) {
	switch txType {
	case TxDeposit, TxWithdrawal:
		return TransactionType(txType), nil
	}

	return "", ErrInvalidValue{Msg: "invalid transaction type"}
}

type Transaction struct {
	ID        string          `json:"id"`
	AccountID string          `json:"accountId"`
	Type      TransactionType `json:"type"`
	Amount    float32         `json:"amount"`
	Timestamp time.Time       `json:"timestamp"`
}

func NewTransaction(id, accountID, txType string, amount float32, timestamp time.Time) (Transaction, error) {
	transactionType, err := NewTransactionType(txType)
	if err != nil {
		return Transaction{}, err
	}

	return Transaction{
		ID:        id,
		AccountID: accountID,
		Type:      transactionType,
		Amount:    amount,
		Timestamp: timestamp,
	}, nil
}
