package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/jyisus/bank-server/internal"
	"github.com/jyisus/bank-server/internal/service"
)

func addRoutes(
	mux *http.ServeMux,
	accountsService *service.AccountService,
	transactionsService *service.TransactionService,
) {
	mux.HandleFunc("POST /accounts", createNewAccountHandler(accountsService))
	mux.HandleFunc("GET /accounts/{id}", retrieveAccountDetails(accountsService))
	mux.HandleFunc("GET /accounts", retrieveAllAccounts(accountsService))
	mux.HandleFunc("POST /accounts/{id}/transactions", createTransactionHandler(transactionsService))
	mux.HandleFunc("GET /accounts/{id}/transactions", retrieveAllTransactions(transactionsService))
	mux.HandleFunc("POST /transfer", transferBetweenAccounts(accountsService))
}

var accountRepo = map[string]internal.Account{}
var transactionsRepo = map[string]internal.Transaction{}

func createNewAccountHandler(
	accountsService *service.AccountService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type createAccountRequest struct {
			Owner          string  `json:"owner"`
			InitialBalance float32 `json:"initial_balance"`
		}

		req, err := decode[createAccountRequest](r)
		if err != nil {
			processError(w, err)
			return
		}

		newAccount, err := accountsService.CreateAccount(context.Background(), req.Owner, req.InitialBalance)
		if err != nil {
			processError(w, err)
			return
		}

		response := struct {
			ID string `json:"id"`
		}{
			ID: newAccount.ID,
		}

		encode(w, http.StatusCreated, response)
	}
}

func retrieveAccountDetails(accountsService *service.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID := r.PathValue("id")

		service, err := accountsService.GetAccount(context.Background(), accountID)
		if err != nil {
			processError(w, err)
			return
		}

		encode(w, http.StatusOK, service)
	}
}

func retrieveAllAccounts(accountsService *service.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {

		accounts, err := accountsService.ListAccounts(context.Background())
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if len(accounts) == 0 {
			encode(w, http.StatusNoContent, accounts)
			return
		}

		encode(w, http.StatusOK, accounts)
	}
}

type CreateTransactionRequest struct {
	Type   string  `json:"type"`
	Amount float32 `json:"amount"`
}

func createTransactionHandler(transactionService *service.TransactionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID := r.PathValue("id")

		transaction, err := decode[CreateTransactionRequest](r)
		if err != nil {
			// slog.Error("Error during body unmarshall", "error", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		tx, err := transactionService.SaveTransaction(
			context.Background(),
			accountID,
			transaction.Type,
			transaction.Amount,
		)
		if err != nil {
			processError(w, err)
			return
		}

		encode(w, http.StatusOK, tx)
	}
}

func retrieveAllTransactions(transactionService *service.TransactionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID := r.PathValue("id")

		transactions, err := transactionService.RetrieveAccountTransactions(context.Background(), accountID)
		if err != nil {
			processError(w, err)
			return
		}

		if len(transactions) == 0 {
			encode(w, http.StatusNoContent, transactions)
			return
		}
		// slog.Info("Retrieving all accounts", "accounts", len(transactions))

		encode(w, http.StatusOK, transactions)
	}
}

func transferBetweenAccounts(accountsService *service.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type tranferRequest struct {
			FromAccountID string  `json:"from_account_id"`
			ToAccountID   string  `json:"to_account_id"`
			Amount        float32 `json:"amount"`
		}

		req, err := decode[tranferRequest](r)
		if err != nil {
			http.Error(w, "Error decoding request", http.StatusBadRequest)
		}

		if err := accountsService.Transfer(
			context.Background(),
			req.FromAccountID,
			req.ToAccountID,
			req.Amount,
		); err != nil {
			processError(w, err)
			return
		}

		encode(w, http.StatusOK, req)
	}
}

func processError(w http.ResponseWriter, err error) {
	switch {
	case errors.As(err, &internal.ErrAccountNotFound{}):
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case errors.As(err, &internal.ErrInsufficientBalance{}):
		// NOTE: I'm using 403 here beacuse it common in this context, but I'm not sure if it's the best fit
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	default:
		slog.Error("Internal server error", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
