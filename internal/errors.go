package internal

import (
	"errors"
	"fmt"
)

// var (
// 	ErrAccountNotFound = errors.New("account not found")
// )

type ErrInvalidValue struct {
	Msg string
}

func (e ErrInvalidValue) Error() string {
	return fmt.Sprintf("invalid value: %s", e.Msg)
}

type ErrAccountNotFound struct {
	AccountID string
}

func (e ErrAccountNotFound) Error() string {
	return fmt.Sprintf("account with id %q not found", e.AccountID)
}

type ErrInsufficientBalance struct {
	AccountID string
}

func (e ErrInsufficientBalance) Error() string {
	return fmt.Sprintf("balance for account with id %q is insufficient", e.AccountID)
}

var (
	ErrAccountAlreadyExists = errors.New("account already exists")
)
