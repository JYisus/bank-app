package internal

import "time"

// TODO: use int?
type TransactionType string

const (
	TxDeposit    = "deposit"
	TxWithdrawal = "withdrawal"
)

type Transaction struct {
	ID        string          `json:"id"`
	AccountID string          `json:"accountId"`
	Type      TransactionType `json:"type"`
	Amount    float32         `json:"amount"`
	Timestamp time.Time       `json:"timestamp"`
}

func NewTransaction(id, accountID, txType string, amount float32, timestamp time.Time) (Transaction, error) {
	return Transaction{
		ID:        id,
		AccountID: accountID,
		Type:      TransactionType(txType),
		Amount:    amount,
		Timestamp: timestamp,
	}, nil
}
