package dto

import (
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
	Balance int `json:"balance"`
}

type CreateAccountInputDTO struct {
	Name    string `json:"name"`
	CPF     string `json:"cpf"`
	Secret  string `json:"secret"`
	Balance int    `json:"balance"`
}

type LoginInputDTO struct {
	CPF    string `json:"cpf"`
	Secret string `json:"secret"`
}

type LoginOutputDTO struct {
	Token string `json:"token"`
}
