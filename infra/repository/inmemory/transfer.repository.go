package inmemory

import (
	"fmt"
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

	fmt.Print(accountOriginID, destinationOriginID, amount)

	transfer := entity.Transfer{
		ID:                   len(r.Transfers),
		AccountOriginID:      accountOriginID,
		AccountDestinationID: destinationOriginID,
		Amount:               amount,
		CreatedAt:            time.Now().UTC(),
	}

	// Falta verificar se é valido tbm
	_, err := transfer.IsValid()
	if err != nil {
		return nil
	}

	r.Transfers = append(r.Transfers, transfer)

	return &transfer

}

func (r *TransferRepository) Reset() error {
	r.Transfers = []entity.Transfer{}
	return nil
}
