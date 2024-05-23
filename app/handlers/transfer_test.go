package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/PPAKruNN/golearn/domain/service"
	"github.com/PPAKruNN/golearn/domain/service/dto"
)

func createHTTPTransferServer() (TransferService *service.TransferService, AccountService *service.AccountService, AuthService *service.AuthService, server *TransferServer) {

	TransferService, AccountService, AuthService = createRepoAndServices()
	server = NewTransferServer(*TransferService, *AuthService)

	return
}

func TestGETTransfers(t *testing.T) {

	TransferService, AccountService, AuthService, server := createHTTPTransferServer()

	createMockAccount(AccountService)
	createMockAccount(AccountService)

	t.Run("Should return Unauthorized if no token is provided!", func(t *testing.T) {

		request, response := createHttpRequestAndResponse(http.MethodGet, "/transfers", nil)
		server.ReadTransfers(response, request)

		assertStatusCode(t, response, http.StatusUnauthorized)
	})

	t.Run("Should return all transfers from the account", func(t *testing.T) {

		request, response := createHttpRequestAndResponse(http.MethodGet, "/transfers", nil)

		// Register a token for account
		token := AuthService.CreateToken(0)
		bearer := "Bearer " + token
		request.Header.Add("Authorization", bearer)

		// Register a transfer
		newTransfer := dto.CreateTrasnferInputDTO{
			AccountOriginID:      0,
			AccountDestinationID: 1,
			Amount:               10,
		}
		_, err := TransferService.CreateTransfer(newTransfer)
		if err != nil {
			t.Errorf("Error while creating mock transfer!")
		}

		server.ReadTransfers(response, request)

		var output []dto.ReadTransfersOutputDTO
		json.NewDecoder(response.Body).Decode(&output)

		assertStatusCode(t, response, http.StatusOK)

		// assertTransfers
		currTransfer := output[0]

		if currTransfer.AccountDestinationID != newTransfer.AccountDestinationID ||
			currTransfer.AccountOriginID != newTransfer.AccountOriginID ||
			currTransfer.Amount != newTransfer.Amount {
			t.Errorf("Persisted transfer is different from sent transfer! \nPersisted: %+v \nSent: %+v", currTransfer, newTransfer)
		}
	})
}

func TestPOSTTransfer(t *testing.T) {
	TransferService, AccountService, AuthService, server := createHTTPTransferServer()

	mockedAccount := createMockAccount(AccountService)
	createMockAccount(AccountService)

	t.Run("Should return Unauthorized if no token is provided!", func(t *testing.T) {

		body, err := json.Marshal(dto.CreateTrasnferInputDTO{
			AccountOriginID:      0,
			AccountDestinationID: 0,
			Amount:               10,
		})
		if err != nil {
			t.Error("Error while creating body for CreateTrasnferInputDTO")
		}

		request, response := createHttpRequestAndResponse(http.MethodPost, "/transfers", bytes.NewBuffer(body))
		server.CreateTransfer(response, request)

		assertStatusCode(t, response, http.StatusUnauthorized)
	})

	t.Run("Should be able to create a transfer", func(t *testing.T) {

		newTransfer := dto.CreateTrasnferInputDTO{
			AccountOriginID:      0,
			AccountDestinationID: 1,
			Amount:               10,
		}

		body, err := json.Marshal(newTransfer)
		if err != nil {
			t.Error("Error while creating body for CreateTrasnferInputDTO")
		}

		request, response := createHttpRequestAndResponse(http.MethodPost, "/transfers", bytes.NewBuffer(body))

		// Register a token for account
		token := AuthService.CreateToken(0)
		bearer := "Bearer " + token
		request.Header.Add("Authorization", bearer)

		oldBalanceOrigin := AccountService.ReadAccountBalance(dto.ReadAccountBalanceInputDTO{ID: 0})
		oldBalanceDest := AccountService.ReadAccountBalance(dto.ReadAccountBalanceInputDTO{ID: 1})

		server.CreateTransfer(response, request)

		newerBalanceOrigin := AccountService.ReadAccountBalance(dto.ReadAccountBalanceInputDTO{ID: 0})
		newerBalanceDest := AccountService.ReadAccountBalance(dto.ReadAccountBalanceInputDTO{ID: 1})

		tranfers := TransferService.ReadTransfersByAccount(0)
		currTransfer := tranfers[0]

		if currTransfer.AccountDestinationID != newTransfer.AccountDestinationID ||
			currTransfer.AccountOriginID != newTransfer.AccountOriginID ||
			currTransfer.Amount != newTransfer.Amount {
			t.Errorf("Persisted transfer is different from sent transfer! \nPersisted: %+v \nSent: %+v", currTransfer, newTransfer)
		}

		assertStatusCode(t, response, http.StatusCreated)

		if (newerBalanceDest.Balance - oldBalanceDest.Balance) !=
			(oldBalanceOrigin.Balance - newerBalanceOrigin.Balance) {
			t.Errorf("Amount transfered is different from received!")
		}

	})

	t.Run("Should NOT be able to create a transfer when account hava insufficient funds", func(t *testing.T) {

		newTransfer := dto.CreateTrasnferInputDTO{
			AccountOriginID:      0,
			AccountDestinationID: 1,
			Amount:               mockedAccount.Balance + 100,
		}

		body, err := json.Marshal(newTransfer)
		if err != nil {
			t.Error("Error while creating body for CreateTrasnferInputDTO")
		}

		request, response := createHttpRequestAndResponse(http.MethodPost, "/transfers", bytes.NewBuffer(body))

		// Register a token for account
		token := AuthService.CreateToken(0)
		bearer := "Bearer " + token
		request.Header.Add("Authorization", bearer)

		server.CreateTransfer(response, request)

		assertStatusCode(t, response, http.StatusBadRequest)

		acc := AccountService.ReadAccountBalance(dto.ReadAccountBalanceInputDTO{ID: 0})
		if acc.Balance < 0 {
			t.Errorf("Transaction removed money from account! It was expected to not do it.")
		}
	})

	t.Run("Should NOT be able to create a transfer to itself", func(t *testing.T) {

		newTransfer := dto.CreateTrasnferInputDTO{
			AccountOriginID:      0,
			AccountDestinationID: 0,
			Amount:               10,
		}

		body, err := json.Marshal(newTransfer)
		if err != nil {
			t.Error("Error while creating body for CreateTrasnferInputDTO")
		}

		request, response := createHttpRequestAndResponse(http.MethodPost, "/transfers", bytes.NewBuffer(body))

		// Register a token for account
		token := AuthService.CreateToken(0)
		bearer := "Bearer " + token
		request.Header.Add("Authorization", bearer)

		server.CreateTransfer(response, request)

		assertStatusCode(t, response, http.StatusBadRequest)
	})

}
