package dto

import "time"

type ReadTransfersOutputDTO struct {
	ID                   int       `json:"id"`
	AccountOriginID      int       `json:"account_origin_id"`
	AccountDestinationID int       `json:"account_destination_id"`
	Amount               int       `json:"amount"`
	CreatedAt            time.Time `json:"created_at"`
}

type CreateTrasnferInputDTO struct {
	AccountOriginID      int `json:"account_origin_id"`
	AccountDestinationID int `json:"account_destination_id"`
	Amount               int `json:"amount"`
}
