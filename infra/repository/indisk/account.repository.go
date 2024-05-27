package indisk

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"os"
	"time"

	"github.com/PPAKruNN/golearn/domain/entity"
)

const (
	ACCOUNT_DATA_FILENAME = "accounts.json"
)

type accountJSONSchema struct {
	ID        int
	Name      string
	CPF       string
	Secret    string
	Balance   int
	CreatedAt time.Time
}

type AccountRepository struct {
	Accounts     []entity.Account
	AccountsJson []accountJSONSchema
	pathToFile   string
}

func NewAccountRepository(dir string) *AccountRepository {

	path := fmt.Sprintf("%s/%s", dir, ACCOUNT_DATA_FILENAME)

	repo := &AccountRepository{
		Accounts:   []entity.Account{{}},
		pathToFile: path,
	}

	repo.loadIntoMemory()

	return repo

}

func (r *AccountRepository) openHandle() *os.File {
	handle := openOrCreateFile(r.pathToFile)
	return handle
}

func (r *AccountRepository) Reset() error {
	err := resetFile(r.pathToFile)
	return err
}

func (r *AccountRepository) save() {
	marshal, _ := json.MarshalIndent(r.AccountsJson, "", "  ")
	saveInFile(r.pathToFile, marshal)
}

func (r *AccountRepository) loadIntoMemory() {
	// Load accounts from disk
	handle := r.openHandle()
	defer handle.Close()

	var accounts []accountJSONSchema

	json.NewDecoder(handle).Decode(&accounts)

	// if err != nil {
	// 	panic("Could not load accounts from disk: ")
	// }

	r.AccountsJson = accounts
	r.Accounts = convertJSONtoEntity(accounts)
}

func convertJSONtoEntity(json []accountJSONSchema) (output []entity.Account) {

	for _, acc := range json {

		hash := sha256.New()
		hash.Write([]byte(acc.Secret))

		output = append(output, entity.Account{
			ID:        acc.ID,
			Name:      acc.Name,
			CPF:       acc.CPF,
			Secret:    hash,
			Balance:   acc.Balance,
			CreatedAt: acc.CreatedAt,
		})
	}

	return
}

func convertEntityToJSON(entities []entity.Account) (output []accountJSONSchema) {

	for _, acc := range entities {

		hash := hex.EncodeToString(acc.Secret.Sum(nil))
		fmt.Println(hash)

		output = append(output, accountJSONSchema{
			ID:        acc.ID,
			Name:      acc.Name,
			CPF:       acc.CPF,
			Secret:    hash,
			Balance:   acc.Balance,
			CreatedAt: acc.CreatedAt,
		})
	}

	return
}

func (r *AccountRepository) ReadAll() []entity.Account {

	r.loadIntoMemory()

	return r.Accounts
}

func (r *AccountRepository) ReadByID(id int) *entity.Account {

	r.loadIntoMemory()

	for idx, acc := range r.Accounts {
		if acc.ID == id {
			return &r.Accounts[idx]
		}
	}

	return nil
}

func (r *AccountRepository) ReadHashByCPF(cpf string) (int, []byte) {

	handle := r.openHandle()
	defer handle.Close()

	var accounts []accountJSONSchema
	json.NewDecoder(handle).Decode(&accounts)

	for _, acc := range accounts {
		if acc.CPF == cpf {

			bs, err := hex.DecodeString(acc.Secret)
			if err != nil {
				panic("Corrupted secret hash")
			}

			return acc.ID, bs
		}
	}

	return 0, []byte{}

}

func (r *AccountRepository) Create(name string, cpf string, secret hash.Hash, balance int) entity.Account {

	r.loadIntoMemory()

	newAccount := entity.Account{
		ID:        len(r.Accounts),
		Name:      name,
		CPF:       cpf,
		Secret:    secret,
		Balance:   balance,
		CreatedAt: time.Now().UTC(),
	}

	r.Accounts = append(r.Accounts, newAccount)
	r.AccountsJson = append(r.AccountsJson, accountJSONSchema{
		ID:        newAccount.ID,
		Name:      newAccount.Name,
		CPF:       newAccount.CPF,
		Secret:    hex.EncodeToString(newAccount.Secret.Sum(nil)),
		Balance:   newAccount.Balance,
		CreatedAt: time.Now().UTC(),
	})

	r.save()

	return newAccount

}

func (r *AccountRepository) UpdateBalance(id, balance int) *entity.Account {

	r.loadIntoMemory()

	for idx, acc := range r.AccountsJson {
		if acc.ID == id {
			r.AccountsJson[idx].Balance = balance
			r.save()

			return &r.Accounts[idx]
		}
	}

	return nil
}
