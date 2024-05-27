package service

import (
	"fmt"
	"net/http"

	"github.com/PPAKruNN/golearn/domain/entity"
	"github.com/PPAKruNN/golearn/domain/service/dto"
)

type TransferRepository interface {
	ReadTransfersByAccountID(id int) []entity.Transfer
	CreateTransfer(accountOriginID, destinationOriginID, amount int) *entity.Transfer
	Reset() error
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

func (t *TransferService) CreateTransfer(input dto.CreateTrasnferInputDTO) (int, error) {

	// FIXME: Remove this 2 queries and turn it into one.
	origin := t.AccountRepo.ReadByID(input.AccountOriginID)
	if origin == nil {
		return http.StatusNotFound, fmt.Errorf("Could not find the origin account!")
	}
	destination := t.AccountRepo.ReadByID(input.AccountDestinationID)
	if destination == nil {
		return http.StatusNotFound, fmt.Errorf("Could not find the destination account!")
	}

	transfer, err := origin.TransferTo(destination, input.Amount)
	if err != nil {
		return http.StatusBadRequest, err
	}

	t.AccountRepo.UpdateBalance(origin.ID, origin.Balance)
	t.AccountRepo.UpdateBalance(destination.ID, destination.Balance)
	persistedTransfer := t.TransferRepo.CreateTransfer(transfer.AccountOriginID, transfer.AccountDestinationID, transfer.Amount)

	// FIXME: Create better error message.
	if persistedTransfer == nil {
		return http.StatusInternalServerError, fmt.Errorf("Error while creating transfer. Internal server error")
	}

	return http.StatusCreated, nil

}
