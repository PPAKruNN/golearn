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
	ReadAll() []entity.Account
	ReadByID(id int) *entity.Account
	ReadByCPF(cpf string) *entity.Account
	Create(name string, cpf string, secret hash.Hash, balance int) entity.Account
}

type AccountService struct {
	Repo     AccountRepository
	AuthRepo AuthRepository
}

func NewAccountService(repo AccountRepository, authRepository AuthRepository) *AccountService {
	return &AccountService{Repo: repo, AuthRepo: authRepository}
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
	account := a.Repo.ReadByID(input.ID)

	if account == nil {
		return dto.ReadAccountBalanceOutputDTO{Balance: -1}
	}

	return dto.ReadAccountBalanceOutputDTO{Balance: account.Balance}
}

func (a AccountService) CreateAccount(input dto.CreateAccountInputDTO) error {

	hash := hashSecret(input.Secret)

	account := entity.NewAccount(
		-1,
		input.Balance,
		input.Name,
		input.CPF,
		hash,
		time.Now().UTC(),
	)

	valid, err := account.IsValid()
	if !valid {
		return err
	}

	_ = a.Repo.Create(input.Name, input.CPF, hash, input.Balance)

	return nil
}

func (a AccountService) Authenticate(cpf, secret string) (int, error) {

	account := a.Repo.ReadByCPF(cpf)
	if account == nil {
		return 0, fmt.Errorf("Failed to authenticate. Cannot find an account with the credentials provided")
	}

	// Checking secrets.
	isCorrectSecret := checkSecret(secret, account.Secret)
	if !isCorrectSecret {
		return 0, fmt.Errorf("Failed to authenticate. Invalid secret provided!")
	}

	return account.ID, nil

}

func hashSecret(secret string) (hash hash.Hash) {
	hash = sha256.New()
	hash.Write([]byte(secret))

	return
}

func checkSecret(secret string, hash hash.Hash) bool {
	newHash := hashSecret(secret)
	return bytes.Equal(newHash.Sum(nil), hash.Sum(nil))
}
