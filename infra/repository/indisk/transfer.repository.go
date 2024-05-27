package indisk

import (
	"encoding/json"
	"os"
	"path"
	"time"

	"github.com/PPAKruNN/golearn/domain/entity"
)

const (
	TRANSFER_DATA_FILENAME = "transfers.json"
)

type TransferRepository struct {
	Transfers  []entity.Transfer
	pathToFile string
}

func NewTransferRepository(dir string) *TransferRepository {
	repo := &TransferRepository{
		Transfers:  []entity.Transfer{},
		pathToFile: path.Join(dir, TRANSFER_DATA_FILENAME),
	}

	repo.loadIntoMemory()

	return repo
}

func (r *TransferRepository) ReadTransfersByAccountID(id int) []entity.Transfer {

	r.loadIntoMemory()

	var transfers []entity.Transfer

	for _, trans := range r.Transfers {
		if trans.AccountOriginID == id {
			transfers = append(transfers, trans)
		}
	}

	return transfers

}

func (r *TransferRepository) CreateTransfer(accountOriginID, destinationOriginID, amount int) *entity.Transfer {

	r.loadIntoMemory()

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

	r.save()

	return &transfer
}

func (r *TransferRepository) openHandle() *os.File {
	handle := openOrCreateFile(r.pathToFile)
	return handle
}

func (r *TransferRepository) Reset() error {
	err := resetFile(r.pathToFile)
	return err
}

func (r *TransferRepository) save() {
	marshal, _ := json.MarshalIndent(r.Transfers, "", "  ")
	saveInFile(r.pathToFile, marshal)
}

func (r *TransferRepository) loadIntoMemory() {
	// Load accounts from disk
	handle := r.openHandle()
	defer handle.Close()

	var transfers []entity.Transfer
	json.NewDecoder(handle).Decode(&transfers)

	r.Transfers = transfers
}
