package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PPAKruNN/golearn/domain/service"
	"github.com/PPAKruNN/golearn/domain/service/dto"
)

type AccountServer struct {
	AccountService  service.AccountService
	TransferService service.TransferService
	AuthService     service.AuthService
}

func NewAccountServer(transferService service.TransferService, accountService service.AccountService, authService service.AuthService) *AccountServer {

	return &AccountServer{
		AccountService:  accountService,
		TransferService: transferService,
		AuthService:     authService,
	}
}

func (s *AccountServer) ServeHTTP() *http.ServeMux {

	router := http.NewServeMux()
	router.Handle("/accounts/{id}/balance", http.HandlerFunc(s.ReadAccountBalance))
	router.Handle("/accounts/", http.HandlerFunc(s.accountHandler))

	return router

}

// FIX: Change this name to something more meaningfull.
func (s *AccountServer) accountHandler(w http.ResponseWriter, r *http.Request) {

	method := r.Method

	switch method {
	case http.MethodPost:
		s.CreateAccount(w, r)
		break
	case http.MethodGet:
		s.ReadAccounts(w, r)
		break

	case "":
		s.ReadAccounts(w, r)
		break

	default:
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Cannot " + r.Method + " " + r.URL.String())
	}

}

func (s *AccountServer) ReadAccounts(w http.ResponseWriter, r *http.Request) {
	accounts := s.AccountService.ReadAccounts()

	json.NewEncoder(w).Encode(accounts)
}

func (s *AccountServer) ReadAccountBalance(w http.ResponseWriter, r *http.Request) {

	// Getting ID
	var id int
	fmt.Sscanf(r.URL.Path, "/accounts/%d/balance", &id)

	input := dto.ReadAccountBalanceInputDTO{
		ID: id,
	}

	balance := s.AccountService.ReadAccountBalance(input)

	if balance.Balance == -1 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(balance)
	}

}

func (s *AccountServer) CreateAccount(w http.ResponseWriter, r *http.Request) {

	var accountDTO dto.CreateAccountInputDTO

	err := json.NewDecoder(r.Body).Decode(&accountDTO)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = s.AccountService.CreateAccount(accountDTO)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Error while creating account! %+v", accountDTO)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (s *AccountServer) Login(w http.ResponseWriter, r *http.Request) {

	var accountDTO dto.LoginInputDTO

	err := json.NewDecoder(r.Body).Decode(&accountDTO)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	id, err := s.AccountService.Authenticate(accountDTO.CPF, accountDTO.Secret)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Printf("Account not found! %+v\n", accountDTO)
		return
	}

	token := s.AuthService.CreateToken(id)

	output := dto.LoginOutputDTO{
		Token: token,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)

}
