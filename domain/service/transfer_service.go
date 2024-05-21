package service

import (
	"fmt"

	"github.com/PPAKruNN/golearn/domain/entity"
	"github.com/PPAKruNN/golearn/domain/service/dto"
)

type TransferRepository interface {
	ReadTransfersByAccountID(id int) []entity.Transfer
	CreateTransfer(accountOriginID, destinationOriginID, amount int) *entity.Transfer
}

type TransferService struct {
	TransferRepo TransferRepository
	AccountRepo  AccountRepository
}

func NewTransferService(transferRepo TransferRepository, accountRepo AccountRepository) *TransferService {
	return &TransferService{TransferRepo: transferRepo, AccountRepo: accountRepo}
}

func (t TransferService) ReadTransfersByAccount(accountId int) []dto.ReadTransfersOutputDTO {

	transfers := t.TransferRepo.ReadTransfersByAccountID(accountId)

	parsedTransfer := []dto.ReadTransfersOutputDTO{}

	// Mapping entity.Transfer to dto.ReadTrasnferOutputDTO
	for _, val := range transfers {

		newTransfer := dto.ReadTransfersOutputDTO{
			ID:                   val.ID,
			AccountOriginID:      val.AccountOriginID,
			AccountDestinationID: val.AccountDestinationID,
			Amount:               val.Amount,
			CreatedAt:            val.CreatedAt,
		}
		parsedTransfer = append(parsedTransfer, newTransfer)

	}

	return parsedTransfer
}

func (t *TransferService) CreateTransfer(input dto.CreateTrasnferInputDTO) error {

	// FIXME: Remove this 2 queries and turn it into one.
	origin := t.AccountRepo.ReadByID(input.AccountOriginID)
	if origin == nil {
		return fmt.Errorf("Could not find the origin account!")
	}
	destination := t.AccountRepo.ReadByID(input.AccountDestinationID)
	if destination == nil {
		return fmt.Errorf("Could not find the destination account!")
	}

	transfer, err := origin.TransferTo(destination, input.Amount)
	if err != nil {
		return err
	}

	persistedTransfer := t.TransferRepo.CreateTransfer(transfer.AccountDestinationID, transfer.AccountDestinationID, transfer.Amount)

	// FIXME: Create better error message.
	if persistedTransfer == nil {
		return fmt.Errorf("Error while creating transfer. Internal server error")
	}

	return nil

}
