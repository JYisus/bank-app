package server

import (
	"net/http"

	"github.com/jyisus/bank-server/internal/service"
)

func New(
	accountsService *service.AccountService,
	transactionsService *service.TransactionService,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, accountsService, transactionsService)

	return mux
}
