package internal

import "regexp"

type Name string

func NewName(name string) (Name, error) {
	// Basic check to avoid names with numbers
	if regexp.MustCompile(`\d`).MatchString(name) {
		return "", ErrInvalidValue{Msg: "owner name can't contain numbers"}
	}

	return Name(name), nil
}

type Account struct {
	ID      string  `json:"id"`
	Owner   Name    `json:"owner"`
	Balance float32 `json:"balance"`
}

func NewAccount(id, owner string, balance float32) (Account, error) {
	if balance < 0 {
		return Account{}, ErrInvalidValue{Msg: "the initial balance should be greater than 0"}
	}

	ownerName, err := NewName(owner)
	if err != nil {
		return Account{}, err
	}

	return Account{
		ID:      id,
		Owner:   ownerName,
		Balance: balance,
	}, nil
}

func (a *Account) Deposit(amount float32) {
	a.Balance += amount
}

func (a *Account) Withdraw(amount float32) error {
	if a.Balance-amount < 0 {
		return ErrInsufficientBalance{AccountID: a.ID}
	}

	a.Balance -= amount

	return nil
}
