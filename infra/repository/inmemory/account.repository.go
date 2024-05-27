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

func (r AccountRepository) UpdateBalance(id int, balance int) *entity.Account {
	for idx, acc := range r.Accounts {
		if acc.ID == id {
			r.Accounts[idx].Balance = balance
			return &r.Accounts[idx]
		}
	}

	return nil
}

func (r AccountRepository) ReadHashByCPF(cpf string) (int, []byte) {

	for _, acc := range r.Accounts {
		if acc.CPF == cpf {
			return acc.ID, acc.Secret.Sum(nil)
		}
	}

	return 0, []byte{}
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

func (r *AccountRepository) Reset() error {
	r.Accounts = []entity.Account{}
	return nil
}
