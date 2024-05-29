package database

import (
	"time"

	"github.com/PPAKruNN/golearn/domain/entity"
	"github.com/jackc/pgx"
	"github.com/rs/zerolog/log"
)

type TransferRepository struct {
	connection *pgx.ConnPool
}

func NewTransferRepository() *TransferRepository {

	connConfig := loadDatabaseEnvs()

	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: connConfig,
	})

	if err != nil {
		log.Error().Err(err).Msg("Unable to connect to database")
		panic("Couldn't connect to database")
	}

	return &TransferRepository{
		connection: pool,
	}

}

func (r *TransferRepository) ReadTransfersByAccountID(id int) []entity.Transfer {

	rows, err := r.connection.Query(`SELECT * FROM "Transfer"`)
	if err != nil {
		log.Info().Err(err).Int("id", id).Msg("Failed to query all transfers from account")
		return nil
	}
	defer rows.Close()

	transfers := []entity.Transfer{}
	for rows.Next() {
		var id int
		var accountOriginID int
		var accountDestinationID int
		var amount int
		var created_at time.Time

		err = rows.Scan(&id, &accountOriginID, &accountDestinationID, &amount, &created_at)
		if err != nil {
			log.Info().Err(err).Msg("Failed to scan transfer")
			return nil
		}

		transfers = append(transfers, entity.Transfer{
			ID:                   id,
			AccountOriginID:      accountOriginID,
			AccountDestinationID: accountDestinationID,
			Amount:               amount,
			CreatedAt:            created_at,
		})
	}

	return transfers
}

func (r *TransferRepository) CreateTransfer(accountOriginID, accountDestinationID, amount int) *entity.Transfer {

	rows, err := r.connection.Query(`INSERT into "Transfer" (account_origin_id, account_destination_id, amount) VALUES ($1, $2, $3) RETURNING id, account_origin_id, account_destination_id, amount, created_at`, accountOriginID, accountDestinationID, amount)

	if err != nil {
		log.Info().Err(err).
			Int("OriginID", accountOriginID).
			Int("DestinationID", accountDestinationID).
			Int("Amount", amount).
			Msg("Failed to create transfer")
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		transfer := entity.Transfer{}
		var id int
		var accountOriginID int
		var accountDestinationID int
		var amount int
		var created_at time.Time

		err = rows.Scan(&id, &accountOriginID, &accountDestinationID, &amount, &created_at)
		if err != nil {
			log.Info().Err(err).Msg("Failed to scan transfer")
			return nil
		}

		transfer = entity.Transfer{
			ID:                   id,
			AccountOriginID:      accountOriginID,
			AccountDestinationID: accountDestinationID,
			Amount:               amount,
			CreatedAt:            created_at,
		}

		return &transfer
	}

	return &entity.Transfer{}
}

func (r *TransferRepository) Reset() error {
	rows, err := r.connection.Query(`DELETE FROM "Transfer"`)

	if err != nil {
		log.Info().Err(err).
			Msg("Failed to delete all Transfers")
		return err
	}
	defer rows.Close()

	return nil
}
