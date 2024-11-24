package main

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/jyisus/bank-server/internal/memrepo"
	"github.com/jyisus/bank-server/internal/server"
	"github.com/jyisus/bank-server/internal/service"
)

var _port = "8080"

func run() error {
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(logHandler)
	accountsRepo := memrepo.NewAccountsRepository()
	transactionsRepo := memrepo.NewTransactionsRepository()

	accountsService := service.NewAccountService(logger, accountsRepo)
	transactionsService := service.NewTransactionService(logger, accountsRepo, transactionsRepo)

	s := server.New(accountsService, transactionsService)

	logger.Info("Server running", "port", _port)

	httpServer := &http.Server{
		Addr:    net.JoinHostPort("localhost", _port),
		Handler: s,
	}

	return httpServer.ListenAndServe()
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
