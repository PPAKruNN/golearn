package inmemory

import (
	"time"

	"github.com/PPAKruNN/golearn/domain/entity"
)

type TransferRepository struct {
	Transfers []entity.Transfer
}

func NewTransferRepository() *TransferRepository {
	return &TransferRepository{
		Transfers: []entity.Transfer{},
	}

}

func (r *TransferRepository) ReadTransfersByAccountID(id int) []entity.Transfer {

	var transfers []entity.Transfer

	for _, trans := range r.Transfers {
		if trans.AccountOriginID == id {
			transfers = append(transfers, trans)
		}
	}

	return transfers

}

func (r *TransferRepository) CreateTransfer(accountOriginID, destinationOriginID, amount int) *entity.Transfer {

	transfer := entity.Transfer{
		ID:                   len(r.Transfers),
		AccountOriginID:      accountOriginID,
		AccountDestinationID: destinationOriginID,
		Amount:               amount,
		CreatedAt:            time.Now().UTC(),
	}
	// Falta verificar se é valido tbm
	transfer.IsValid()

	r.Transfers = append(r.Transfers, transfer)

	return &transfer

}
