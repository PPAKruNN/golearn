package database

import (
	"encoding/hex"
	"fmt"
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
		ConnConfig: loadDatabaseEnvs(),
	})

	if err != nil {
		log.Error().Err(err).Msg("Unable to connect to database")
		panic("Couldn't connect to database")
	}

	return &AccountRepository{connection: pool}
}

func (r *AccountRepository) ReadAll() ([]entity.Account, error) {

	rows, err := r.connection.Query(`SELECT * FROM "Account"`)
	if err != nil {
		log.Info().Err(err).Msg("Failed to query all accounts")
		return []entity.Account{}, err
	}
	defer rows.Close()

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
			log.Info().Err(err).Msg("Failed to scan account")
			return accounts, err
		}

		accounts = append(accounts, entity.Account{
			ID:        id,
			CPF:       cpf,
			Name:      name,
			Balance:   balance,
			CreatedAt: created_at,
		})
	}

	return accounts, nil

}

func (r *AccountRepository) ReadByID(id int) (entity.Account, error) {

	rows, err := r.connection.Query(`SELECT * FROM "Account" WHERE id = $1`, id)

	if err != nil {
		log.Info().Err(err).Msg("Failed to query accounts")
		return entity.Account{}, err
	}
	defer rows.Close()

	for rows.Next() {

		account := entity.Account{}

		var id int
		var name string
		var cpf string
		var secret string
		var balance int
		var created_at time.Time

		err = rows.Scan(&id, &name, &cpf, &secret, &balance, &created_at)
		if err != nil {
			log.Info().Err(err).Msg("Failed to scan account")
			return entity.Account{}, err
		}

		account = entity.Account{
			ID:        id,
			Name:      name,
			CPF:       cpf,
			Balance:   balance,
			CreatedAt: created_at,
		}

		return account, nil
	}

	return entity.Account{}, fmt.Errorf("Couldn't find a account with the provided ID")

}

func (r *AccountRepository) ReadHashByCPF(cpf string) (int, []byte, error) {

	rows, err := r.connection.Query(`SELECT id, secret FROM "Account" WHERE cpf = $1`, cpf)
	if err != nil {
		log.Info().Err(err).Str("CPF", cpf).Msg("Failed to query account hash by cpf")
		return 0, []byte{}, err
	}
	defer rows.Close()

	var id int
	var hash string

	for rows.Next() {

		err = rows.Scan(&id, &hash)
		if err != nil {
			log.Info().Err(err).Msg("Failed to scan id and hash of account")
			return 0, []byte{}, err
		}
	}

	bs, err := hex.DecodeString(hash)
	if err != nil {
		log.Info().Err(err).Str("hash", hash).Msg("Failed to decode hash provided")
		return 0, []byte{}, err
	}

	return id, bs, nil
}

func (r *AccountRepository) Create(acc entity.Account) (entity.Account, error) {

	log.Info().Str("name", acc.Name).Str("cpf", acc.CPF).Int("balance", acc.Balance).Msg("Creating account")

	encoded := hex.EncodeToString(acc.Secret.Sum(nil))
	rows, err := r.connection.Query(`INSERT INTO "Account" (name, cpf, secret, balance) VALUES ($1, $2, $3, $4) RETURNING id, name, cpf, secret, balance, created_at`, acc.Name, acc.CPF, encoded, acc.Balance)

	if err != nil {
		log.Info().Err(err).Interface("account", acc).Msg("Failed to create account")
		return entity.Account{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var cpf string
		var secret string
		var balance int
		var created_at time.Time

		err = rows.Scan(&id, &name, &cpf, &secret, &balance, &created_at)
		if err != nil {
			log.Info().Err(err).Msg("Failed to scan account")
			return entity.Account{}, err
		}
		acc.ID = id
	}

	return acc, nil
}

func (r *AccountRepository) UpdateBalance(id, balance int) (entity.Account, error) {

	rows, err := r.connection.Query(`UPDATE "Account" SET balance = $1 WHERE id = $2 RETURNING id, name, cpf, secret, balance, created_at`, balance, id)
	if err != nil {
		log.Info().Err(err).Int("id", id).Int("balance", balance).Msg("Failed to update account balance")
		return entity.Account{}, err
	}
	defer rows.Close()

	for rows.Next() {
		account := entity.Account{}

		var id int
		var name string
		var cpf string
		var secret string
		var balance int
		var created_at time.Time

		err = rows.Scan(&id, &name, &cpf, &secret, &balance, &created_at)
		if err != nil {
			log.Info().Err(err).Msg("Failed to scan account")
			return entity.Account{}, err
		}

		account = entity.Account{
			ID:        id,
			Name:      name,
			CPF:       cpf,
			Balance:   balance,
			CreatedAt: created_at,
		}

		return account, nil
	}

	return entity.Account{}, fmt.Errorf("Couldn't UPDATE account")

}

func (r *AccountRepository) Reset() error {

	_, err := r.connection.Query(`DELETE FROM "Account"`)
	if err != nil {
		log.Info().Err(err).Msg("Failed to reset accounts")
		return err
	}
	return nil

}
