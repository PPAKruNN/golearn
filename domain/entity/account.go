package entity

import (
	"fmt"
	"hash"
	"time"
)

type Account struct {
	ID        int
	Name      string
	CPF       string
	Secret    hash.Hash
	Balance   int
	CreatedAt time.Time
}

func NewAccount(id, balance int, name string, cpf string, secret hash.Hash, createdAt time.Time) *Account {
	return &Account{
		ID:        id,
		Balance:   balance,
		Name:      name,
		CPF:       cpf,
		Secret:    secret,
		CreatedAt: createdAt,
	}
}

func (a Account) IsValid() (bool, error) {

	if a.Balance < 0 {
		return false, fmt.Errorf("Cannot have a account with negative balance! Account ID: %d, Balance: %d", a.ID, a.Balance)
	}

	return true, nil

}

func (a *Account) Transfer(amount int, destination *Account) (Transfer, error) {

	if a.Balance < amount {
		return Transfer{}, fmt.Errorf("Cannot create transfer because insuficiend funds. Account Balance: %d, Transfer amount: %d", a.Balance, amount)
	}

	// Temporary ID data just for understanding
	transfer := NewTransfer(1, a.ID, destination.ID, amount, time.Now())

	valid, err := transfer.IsValid()
	if !valid {
		return Transfer{}, err
	}

	// Transfer digital money.
	a.Balance -= amount
	destination.Balance += amount

	return *transfer, nil
}
