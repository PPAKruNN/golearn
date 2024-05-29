package handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PPAKruNN/golearn/domain/entity"
	"github.com/PPAKruNN/golearn/domain/service"
	"github.com/PPAKruNN/golearn/domain/service/dto"
	"github.com/PPAKruNN/golearn/infra/repository/database"
	"github.com/google/uuid"
)

const (
	MOCKED_SECRET  = "senhaSegura"
	MOCKED_NAME    = "Zé Goopher"
	MOCKED_CPF     = "99988877760"
	MOCKED_BALANCE = 100
)

func createMockAccount(AccountService *service.AccountService) entity.Account {

	secret := sha256.New()
	_, err := secret.Write([]byte(MOCKED_SECRET))
	if err != nil {
		panic("Cannot create a hash of mock account secret!")
	}

	accountEntity := entity.Account{
		Name:    MOCKED_NAME,
		CPF:     MOCKED_CPF,
		Secret:  secret,
		Balance: MOCKED_BALANCE,
	}

	account, err := AccountService.Repo.Create(accountEntity)
	if err != nil {
		panic("Cannot create mock account! Err: " + err.Error())
	}

	return account
}

func createRepoAndServices() (TransferService *service.TransferService, AccountService *service.AccountService, AuthService *service.AuthService) {

	// dir, err := os.Getwd()
	// if err != nil {
	// 	panic("Failed to get current directory.")
	// }

	transferRepo := database.NewTransferRepository()
	authRepo := database.NewAuthRepository()
	accountRepo := database.NewAccountRepository()

	accountRepo.Reset()
	authRepo.Reset()
	transferRepo.Reset()

	TransferService = service.NewTransferService(transferRepo, accountRepo)
	AccountService = service.NewAccountService(accountRepo, authRepo)
	AuthService = service.NewAuthService(authRepo)

	return

}

func createHTTPAccountServer() (TransferService *service.TransferService, AccountService *service.AccountService, AuthService *service.AuthService, server *AccountServer) {

	TransferService, AccountService, AuthService = createRepoAndServices()

	server = NewAccountServer(*TransferService, *AccountService, *AuthService)

	return
}

func assertStatusCode(t *testing.T, response *httptest.ResponseRecorder, expected int) {
	t.Helper()
	got := response.Code

	if got != expected {
		t.Errorf("Expected status %d. But got: status %d \n", expected, got)
	}
}

func assertAccountsFromResponse(t *testing.T, response *httptest.ResponseRecorder, accounts []dto.ReadAccountOutputDTO) {
	t.Helper()

	var got []dto.ReadAccountOutputDTO
	err := json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Errorf("Error while decoding accounts json array!. Got: %+v \n", got)
	}

	for _, acc := range accounts {

		var persistedAccount *dto.ReadAccountOutputDTO

		// Searchs by the account ID within array.
		for _, cmp := range got {
			if cmp.ID == acc.ID {
				persistedAccount = &cmp
			}
		}

		if persistedAccount == nil {
			t.Errorf("Response does not contain an specified account. Searched by: %+v", acc)
			return
		}

		// Assert if equal:
		if acc.Name != persistedAccount.Name ||
			acc.Balance != persistedAccount.Balance ||
			acc.CPF != persistedAccount.CPF {

			t.Errorf("Account returned by endpoint is not equal as the provided. \nGot: %+v, \nExpected: %+v ", persistedAccount, acc)
			return
		}
	}

}

func createHttpRequestAndResponse(method, URL string, body io.Reader) (*http.Request, *httptest.ResponseRecorder) {

	request, _ := http.NewRequest(method, URL, body)
	response := httptest.NewRecorder()

	return request, response
}

func TestGETAccounts(t *testing.T) {

	_, AccountService, _, server := createHTTPAccountServer()

	t.Run("Should return empty array when no account is created", func(t *testing.T) {

		t.Cleanup(func() {
			clearDatabase(server)
		})

		request, response := createHttpRequestAndResponse(http.MethodGet, "/accounts", nil)
		server.ReadAccounts(response, request)

		assertStatusCode(t, response, http.StatusOK)
		assertAccountsFromResponse(t, response, []dto.ReadAccountOutputDTO{})
	})

	t.Run("Should return accounts when a there is any registered account", func(t *testing.T) {

		t.Cleanup(func() {
			clearDatabase(server)
		})

		mockAccount := createMockAccount(AccountService)
		account := dto.ReadAccountOutputDTO{
			ID:      mockAccount.ID,
			Name:    mockAccount.Name,
			CPF:     mockAccount.CPF,
			Balance: mockAccount.Balance,
		}

		request, response := createHttpRequestAndResponse(http.MethodGet, "/accounts", nil)
		server.ReadAccounts(response, request)

		assertStatusCode(t, response, http.StatusOK)
		assertAccountsFromResponse(t, response, []dto.ReadAccountOutputDTO{account})
	})
}

func TestGetBalance(t *testing.T) {

	_, AccountService, _, server := createHTTPAccountServer()

	t.Run("Should return account list", func(t *testing.T) {

		t.Cleanup(func() {
			clearDatabase(server)
		})

		mockAccount := createMockAccount(AccountService)

		fmt.Print(mockAccount.ID)

		request, response := createHttpRequestAndResponse(http.MethodGet, fmt.Sprintf("/accounts/%d/balance", mockAccount.ID), nil)
		server.ReadAccountBalance(response, request)

		assertStatusCode(t, response, http.StatusOK)
		assertAccountBalance(t, response, mockAccount.Balance)
	})

	t.Run("Should return 404 if account is not found!", func(t *testing.T) {

		t.Cleanup(func() {
			clearDatabase(server)
		})

		// FIXME: Should use a number that doesnt exist. Not this hardcoded magic number!
		request, response := createHttpRequestAndResponse(http.MethodGet, fmt.Sprintf("/accounts/%d/balance", 31039813098), nil)
		server.ReadAccountBalance(response, request)

		assertStatusCode(t, response, http.StatusNotFound)
		assertBody(t, response, "")
	})
}

func assertBody(t *testing.T, response *httptest.ResponseRecorder, expectedBody string) {
	t.Helper()
	body := response.Body.String()
	if body != expectedBody {
		t.Errorf("Body is different from expected! Got: %s, expected: %s", body, expectedBody)
	}
}

func assertAccountBalance(t *testing.T, response *httptest.ResponseRecorder, expectedBalance int) {
	t.Helper()

	var got dto.ReadAccountBalanceOutputDTO
	err := json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Errorf("Failed parsing %s to ReadAccountBalanceOutputDTO", response.Body)
		return
	}

	if got.Balance != expectedBalance {
		t.Errorf("Balance is different from expected! Got: %d, expected: %d", got, expectedBalance)
		return
	}

}

// FIX: Refactor this function to become more DRY code.
func TestCreateAccount(t *testing.T) {

	_, AccountService, _, server := createHTTPAccountServer()

	t.Run("Should create a account", func(t *testing.T) {

		t.Cleanup(func() {
			clearDatabase(server)
		})

		// Arrange
		input := dto.CreateAccountInputDTO{
			Name:    MOCKED_NAME,
			CPF:     MOCKED_CPF,
			Secret:  MOCKED_SECRET,
			Balance: MOCKED_BALANCE,
		}
		jsonInput, err := json.Marshal(input)
		if err != nil {
			t.Errorf("Error while converting CreateAccountInputDTO to json. Err: %v", err)
			return
		}

		// Act
		request, response := createHttpRequestAndResponse(http.MethodPost, "/accounts", bytes.NewBuffer(jsonInput))
		server.CreateAccount(response, request)

		persistedAccounts, err := AccountService.Repo.ReadAll()
		if err != nil {
			t.Errorf("Cannot find accounts. Err: %v", err)
			return
		}

		if len(persistedAccounts) != 1 {
			t.Errorf("Account was not created! Expected accounts length to be == 1, but it is: %d", len(persistedAccounts))
			return
		}

		currAccount := persistedAccounts[0]
		if currAccount.CPF != input.CPF ||
			currAccount.Balance != input.Balance ||
			currAccount.Name != input.Name {

			t.Errorf("Account was badly created! Account with ID: 0 does is not equal to input. \nExpected: %+v, \nGot: %+v", input, currAccount)
			return
		}

		assertStatusCode(t, response, http.StatusCreated)
	})

}

// FIX: Refactor this function to more DRY code.
func TestLogin(t *testing.T) {

	_, AccountService, _, server := createHTTPAccountServer()

	t.Run("Should be able to login", func(t *testing.T) {

		t.Cleanup(func() {
			AccountService.Repo.Reset()
		})

		mockedAccount := createMockAccount(AccountService)

		inputDTO := dto.LoginInputDTO{
			CPF:    mockedAccount.CPF,
			Secret: MOCKED_SECRET,
		}

		input, marshalErr := json.Marshal(inputDTO)
		if marshalErr != nil {
			t.Errorf("Failed decoding json into LoginOutputDTO.")
		}

		request, response := createHttpRequestAndResponse(http.MethodPost, "/login", bytes.NewBuffer(input))
		server.Login(response, request)

		var output dto.LoginOutputDTO
		err := json.NewDecoder(response.Body).Decode(&output)
		if err != nil {
			t.Errorf("Failed decoding json into LoginOutputDTO.")
		}

		// Check if token is a valid UUID.
		if uuid.Validate(output.Token) != nil {
			t.Errorf("Login returned string is not a valid UUID.")
		}

	})
}

func clearDatabase(server *AccountServer) {
	server.AccountService.Repo.Reset()
	server.TransferService.TransferRepo.Reset()
	server.AuthService.Repo.Reset()
}
