package database

import (
	"encoding/hex"
	"fmt"
	"hash"
	"os"
	"time"

	"github.com/PPAKruNN/golearn/domain/entity"
	"github.com/jackc/pgx"
	"github.com/rs/zerolog/log"
)

type AccountRepository struct {
	connection *pgx.ConnPool
}

func NewAccountRepository() *AccountRepository {

	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			Port:     5432,
			Database: "golearn",
			User:     "postgres",
			Password: "postgres",
		},
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return &AccountRepository{connection: pool}
}

func (r *AccountRepository) ReadAll() []entity.Account {

	rows, err := r.connection.Query(`SELECT * FROM "Account"`)
	if err != nil {
		log.Error().Err(err).Msg("Failed to query all accounts")
		return nil
	}

	fmt.Printf("\n\n\n\n\n\nReturned rows: %+v\n\n\n\n\n\n\n", rows)

	accounts := []entity.Account{}
	for rows.Next() {
		var id int
		var name string
		var cpf string
		var secret string
		var balance int
		var created_at time.Time

		err = rows.Scan(&id, &name, &cpf, &secret, &balance, &created_at)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan account")
			return nil
		}

		accounts = append(accounts, entity.Account{
			ID:      id,
			Name:    name,
			Balance: balance,
		})
	}

	return accounts

}

func (r *AccountRepository) ReadByID(id int) *entity.Account {

	rows, err := r.connection.Query(`SELECT * FROM "Account" WHERE id = $1`, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to query all accounts")
		return nil
	}

	account := entity.Account{}
	for rows.Next() {
		var id int
		var name string
		var balance int

		err = rows.Scan(&id, &name, &balance)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan account")
			return nil
		}

		account = entity.Account{
			ID:      id,
			Name:    name,
			Balance: balance,
		}
	}

	return &account
}

func (r *AccountRepository) ReadHashByCPF(cpf string) (int, []byte) {

	rows, err := r.connection.Query(`SELECT secret FROM "Account" WHERE cpf = $1`, cpf)
	if err != nil {
		log.Error().Err(err).Str("CPF", cpf).Msg("Failed to query account hash by cpf")
		return 0, []byte{}
	}

	var id int
	var hash string

	for rows.Next() {

		err = rows.Scan(&id, &hash)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan id and hash of account")
			return 0, []byte{}
		}
	}

	bs, err := hex.DecodeString(hash)
	if err != nil {
		log.Error().Err(err).Str("hash", hash).Msg("Failed to decode hash provided")
		return 0, []byte{}
	}

	return id, bs
}

func (r *AccountRepository) Create(name string, cpf string, secret hash.Hash, balance int) entity.Account {

	log.Info().Str("name", name).Str("cpf", cpf).Int("balance", balance).Msg("Creating account")

	encoded := hex.EncodeToString(secret.Sum(nil))

	account := entity.Account{
		Name:    name,
		CPF:     cpf,
		Secret:  secret,
		Balance: balance,
	}

	fmt.Println("Called query with values: ", name, cpf, encoded, balance)

	rows, err := r.connection.Query(`INSERT INTO "Account" (name, cpf, secret, balance) VALUES ($1, $2, $3, $4)`, name, cpf, encoded, balance)

	rows.Close()
	if err != nil {
		log.Error().Err(err).Interface("account", account).Msg("Failed to create account")
		return *&entity.Account{}
	}

	for rows.Next() {
		var id int

		err = rows.Scan(&id)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan account")
			return *&entity.Account{}
		}
		account.ID = id
	}

	return account
}

func (r *AccountRepository) UpdateBalance(id, balance int) *entity.Account {

	_, err := r.connection.Query(`UPDATE "Account" SET balance = $1 WHERE id = $2`, balance, id)
	if err != nil {
		log.Error().Err(err).Int("id", id).Int("balance", balance).Msg("Failed to update account balance")
		return nil
	}

	// FIXME: This is slow, change it to something better.
	account := r.ReadByID(id)
	return account

}

func (r *AccountRepository) Reset() error {

	_, err := r.connection.Query(`DELETE FROM "Account"`)
	if err != nil {
		log.Error().Err(err).Msg("Failed to reset accounts")
		return err
	}
	return nil

}
