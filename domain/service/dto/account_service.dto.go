package dto

import (
	"hash"
	"time"
)

type ReadAccountOutputDTO struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CPF       string    `json:"cpf"`
	Balance   int       `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type ReadAccountBalanceInputDTO struct {
	ID int `json:"id"`
}

type ReadAccountBalanceOutputDTO struct {
	ID int `json:"id"`
}

type CreateAccountInputDTO struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CPF       string    `json:"cpf"`
	Secret    hash.Hash `json:"secret"`
	Balance   int       `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}
