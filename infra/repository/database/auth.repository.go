package database

import (
	"fmt"

	"github.com/jackc/pgx"
	"github.com/rs/zerolog/log"
)

type AuthRepository struct {
	connection *pgx.ConnPool
}

func NewAuthRepository() *AuthRepository {

	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: loadDatabaseEnvs(),
	})

	if err != nil {
		log.Error().Err(err).Msg("Unable to connect to database")
		panic("Couldn't connect to database")
	}

	return &AuthRepository{
		connection: pool,
	}
}

func (r AuthRepository) RegisterToken(token string, accountId int) {
	rows, err := r.connection.Query(`INSERT INTO "Auth" (id, token) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET token = $2`, accountId, token)

	if err != nil {
		log.Info().Err(err).Str("Token", token).Msg("Failed to register token on database!")
		// FIXME: Should return error!
		return
	}

	defer rows.Close()
}

func (r AuthRepository) DecodeToken(token string) (int, error) {

	rows, err := r.connection.Query(`SELECT id FROM "Auth" WHERE token = $1`, token)

	if err != nil {
		log.Info().Err(err).Str("Token", token).Msg("Failed to decode/get accountID from token!")
		return 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var id int

		err = rows.Scan(&id)
		if err != nil {
			log.Info().Err(err).Msg("Failed to scan token")
			return 0, nil
		}

		return id, nil
	}

	return 0, fmt.Errorf("Failed to get decode/get accountId from token!")
}

func (r AuthRepository) Reset() error {
	rows, err := r.connection.Query(`DELETE FROM "Auth"`)

	if err != nil {
		log.Info().Err(err).Msg("Failed reset tokens from database!")
		return err
	}

	defer rows.Close()
	return nil
}
