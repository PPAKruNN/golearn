package entity

import (
	"fmt"
	"time"
)

type Transfer struct {
	ID                   int
	AccountOriginID      int
	AccountDestinationID int
	Amount               int
	CreatedAt            time.Time
}

func NewTransfer(id, accountOriginID, accountDestinationID, amount int, createdAt time.Time) *Transfer {
	return &Transfer{
		ID:                   id,
		AccountOriginID:      accountOriginID,
		AccountDestinationID: accountDestinationID,
		Amount:               amount,
		CreatedAt:            createdAt,
	}
}

func (t Transfer) IsValid() (bool, error) {

	if t.Amount <= 0 {
		return false, fmt.Errorf("Transfer amount cannot be less than 1. Transfer ID: %d, Amount: %d", t.ID, t.Amount)
	}

	if t.AccountDestinationID == t.AccountOriginID {
		return false, fmt.Errorf("Transfer cannot have itself as destination. Transfer ID: %d", t.ID)
	}

	return true, nil
}
