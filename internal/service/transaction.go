package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jyisus/bank-server/internal"
)

type TransactionService struct {
	logger                 *slog.Logger
	accountsRepository     internal.AccountsRepository
	transactionsRepository internal.TransactionsRepository
}

func NewTransactionService(
	logger *slog.Logger,
	accountsRepository internal.AccountsRepository,
	transactionsRepository internal.TransactionsRepository,
) *TransactionService {
	return &TransactionService{
		logger:                 logger,
		accountsRepository:     accountsRepository,
		transactionsRepository: transactionsRepository,
	}
}

func (s TransactionService) SaveTransaction(
	ctx context.Context,
	accountID,
	txType string,
	amount float32,
) (*internal.Transaction, error) {
	txID := uuid.NewString()

	transaction, err := internal.NewTransaction(
		txID,
		accountID,
		txType,
		amount,
		time.Now(),
	)
	if err != nil {
		return nil, err
	}

	account, err := s.accountsRepository.Get(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("getting transaction's account: %w", err)
	}

	switch transaction.Type {
	case internal.TxDeposit:
		account.Deposit(transaction.Amount)
	case internal.TxWithdrawal:
		if err := account.Withdraw(transaction.Amount); err != nil {
			return nil, err
		}
	}

	// if txResult < 0 {
	// 	http.Error(w, "Invalid transaction: not enough money", http.StatusForbidden)
	// 	return
	// }
	newTransaction := internal.Transaction{
		ID:        txID,
		AccountID: accountID,
		Type:      internal.TransactionType(transaction.Type),
		Amount:    transaction.Amount,
		Timestamp: time.Now(),
	}

	if err := s.accountsRepository.UpdateBalance(ctx, accountID, account.Balance); err != nil {
		return nil, fmt.Errorf("updating account balance: %w", err)
	}

	if err := s.transactionsRepository.Save(ctx, newTransaction); err != nil {
		return nil, fmt.Errorf("saving transaction: %w", err)
	}

	return &newTransaction, nil
}

func (s TransactionService) RetrieveAccountTransactions(
	ctx context.Context,
	accountID string,
) ([]internal.Transaction, error) {
	if _, err := s.accountsRepository.Get(ctx, accountID); err != nil {
		return nil, err
	}

	transactions, err := s.transactionsRepository.FindAllByAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
