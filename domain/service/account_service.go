package service

import (
	"hash"
	"time"

	"github.com/PPAKruNN/golearn/domain/entity"
	"github.com/PPAKruNN/golearn/domain/service/dto"
)

type AccountRepository interface {
	ReadAll() []entity.Account
	ReadById(id int) entity.Account
	Create(id int, name string, cpf string, secret hash.Hash, balance int, createdAt time.Time) entity.Account
}

type AccountService struct {
	Repo AccountRepository
}

func NewAccountService(repo AccountRepository) *AccountService {
	return &AccountService{Repo: repo}
}

func (a AccountService) ReadAccounts() []dto.ReadAccountOutputDTO {

	accounts := a.Repo.ReadAll()
	var mapAccounts []dto.ReadAccountOutputDTO

	for _, account := range accounts {

		dto := dto.ReadAccountOutputDTO{
			ID:        account.ID,
			Name:      account.Name,
			CPF:       account.CPF,
			Balance:   account.Balance,
			CreatedAt: account.CreatedAt,
		}

		mapAccounts = append(mapAccounts, dto)
	}

	return mapAccounts
}

func (a AccountService) ReadAccountBalance(input dto.ReadAccountBalanceInputDTO) dto.ReadAccountBalanceOutputDTO {
	account := a.Repo.ReadById(input.ID)

	return dto.ReadAccountBalanceOutputDTO{ID: account.ID}
}

func (a AccountService) CreateAccount(input dto.CreateAccountInputDTO) error {

	account := entity.Account{
		ID:        input.ID,
		Name:      input.Name,
		CPF:       input.CPF,
		Secret:    input.Secret,
		Balance:   input.Balance,
		CreatedAt: input.CreatedAt,
	}

	valid, err := account.IsValid()
	if !valid {
		return err
	}

	_ = a.Repo.Create(input.ID, input.Name, input.CPF, input.Secret, input.Balance, input.CreatedAt)

	return nil
}
