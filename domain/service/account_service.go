package service

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"hash"
	"time"

	"github.com/PPAKruNN/golearn/domain/entity"
	"github.com/PPAKruNN/golearn/domain/service/dto"
)

type AccountRepository interface {
	Create(entity.Account) (entity.Account, error)
	ReadAll() ([]entity.Account, error)
	// FindByID(id int) (entity.Account, error)
	ReadByID(id int) (entity.Account, error)
	// FindHashByCPF(cpf string) (int, []byte, error)
	ReadHashByCPF(cpf string) (int, []byte, error)
	UpdateBalance(id, balance int) (entity.Account, error)
	Reset() error
}

type AccountService struct {
	Repo     AccountRepository
	AuthRepo AuthRepository
}

func NewAccountService(repo AccountRepository, authRepository AuthRepository) *AccountService {
	return &AccountService{Repo: repo, AuthRepo: authRepository}
}

func (a AccountService) ReadAccounts() []dto.ReadAccountOutputDTO {

	accounts, err := a.Repo.ReadAll()
	if err != nil {
	}

	mapAccounts := []dto.ReadAccountOutputDTO{}

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
	account, err := a.Repo.ReadByID(input.ID)

	if err != nil {
		return dto.ReadAccountBalanceOutputDTO{Balance: -1}
	}

	return dto.ReadAccountBalanceOutputDTO{Balance: account.Balance}
}

func (a AccountService) CreateAccount(input dto.CreateAccountInputDTO) (entity.Account, error) {

	hash := hashSecret(input.Secret)

	account := entity.NewAccount(
		0,
		input.Balance,
		input.Name,
		input.CPF,
		hash,
		time.Now(),
	)

	valid, err := account.IsValid()
	if !valid {
		return entity.Account{}, err
	}

	newAccount, err := a.Repo.Create(*account)

	if err != nil {
		return entity.Account{}, err
	}

	return newAccount, nil
}

func (a AccountService) Authenticate(cpf, secret string) (int, error) {

	id, foundSecret, err := a.Repo.ReadHashByCPF(cpf)

	if err != nil {
		return 0, fmt.Errorf("Failed to authenticate. Cannot find an account with the credentials provided")
	}

	// Checking secrets.
	isCorrectSecret := checkSecret(secret, foundSecret)
	if !isCorrectSecret {
		return 0, fmt.Errorf("Failed to authenticate. Invalid secret provided!")
	}

	return id, nil

}

func hashSecret(secret string) (hash hash.Hash) {
	hash = sha256.New()
	hash.Write([]byte(secret))

	return
}

func checkSecret(secret string, accSecret []byte) bool {
	newHash := hashSecret(secret)

	fhash := newHash.Sum(nil)

	return bytes.Equal(fhash, accSecret)
}
