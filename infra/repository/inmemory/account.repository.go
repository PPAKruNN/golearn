package inmemory

import (
	"hash"
	"time"

	"github.com/PPAKruNN/golearn/domain/entity"
)

type AccountRepository struct {
	Accounts []entity.Account
}

func NewAccountRepository() *AccountRepository {
	return &AccountRepository{
		Accounts: []entity.Account{},
	}
}

func (r AccountRepository) ReadAll() []entity.Account {
	return r.Accounts
}

func (r AccountRepository) ReadByID(id int) *entity.Account {
	for idx, acc := range r.Accounts {
		if acc.ID == id {
			return &r.Accounts[idx]
		}
	}

	return nil
}

func (r AccountRepository) ReadByCPF(cpf string) *entity.Account {

	for idx, acc := range r.Accounts {
		if acc.CPF == cpf {
			return &r.Accounts[idx]
		}
	}

	return nil
}

func (r *AccountRepository) Create(name string, cpf string, secret hash.Hash, balance int) entity.Account {

	newAccount := entity.Account{
		ID:        len(r.Accounts),
		Name:      name,
		CPF:       cpf,
		Secret:    secret,
		Balance:   balance,
		CreatedAt: time.Now().UTC(),
	}

	r.Accounts = append(r.Accounts, newAccount)

	return newAccount

}
