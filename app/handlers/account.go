package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PPAKruNN/golearn/domain/service"
	"github.com/PPAKruNN/golearn/domain/service/dto"
	"github.com/rs/zerolog/log"
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

	log.Info().Str("Method", r.Method).Str("Path", r.URL.String()).Msg("Called endpoint ReadAccounts!")

	accounts := s.AccountService.ReadAccounts()

	json.NewEncoder(w).Encode(accounts)

	// INFO: I think it is not worth to put accounts array in the log.
	log.Info().
		Str("Method", r.Method).
		Str("Path", r.URL.String()).
		Int("Status Code", http.StatusOK).
		Msg("")
}

func (s *AccountServer) ReadAccountBalance(w http.ResponseWriter, r *http.Request) {

	log.Info().Str("Method", r.Method).Str("Path", r.URL.String()).Msg("Called endpoint ReadAccountBalance!")

	// Getting ID
	var id int
	fmt.Sscanf(r.URL.Path, "/accounts/%d/balance", &id)

	// Getting ID
	// var id int
	// rawId := r.PathValue("id")
	// parsedInt, err := strconv.Atoi(rawId)
	//
	// if err != nil {
	// 	w.WriteHeader(http.StatusUnprocessableEntity)
	// 	log.Info().
	// 		Err(err).
	// 		Str("Method", r.Method).
	// 		Str("Path", r.URL.String()).
	// 		Int("Status Code", http.StatusUnprocessableEntity).
	// 		Int("id", id).
	// 		Str("rawId", rawId).
	// 		Msg("Couldn't read the id from request")
	//
	// 	return
	// }

	// r.Context()
	input := dto.ReadAccountBalanceInputDTO{
		ID: id,
	}

	balance := s.AccountService.ReadAccountBalance(input)

	if balance.Balance == -1 {
		w.WriteHeader(http.StatusNotFound)
		log.Info().
			Str("Method", r.Method).
			Str("Path", r.URL.String()).
			Int("Status Code", http.StatusNotFound).
			Interface("Account", input).
			Msg("Could not read balance from account!")
	} else {
		json.NewEncoder(w).Encode(balance)

		log.Info().
			Str("Method", r.Method).
			Str("Path", r.URL.String()).
			Int("Status Code", http.StatusOK).
			Interface("Account", input).
			Interface("Response", balance).
			Msg("")
	}

}

func (s *AccountServer) CreateAccount(w http.ResponseWriter, r *http.Request) {

	log.Info().Str("Method", r.Method).Str("Path", r.URL.String()).Msg("Called endpoint CreateAccount!")

	var accountDTO dto.CreateAccountInputDTO

	err := json.NewDecoder(r.Body).Decode(&accountDTO)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Info().
			Str("Method", r.Method).
			Str("Path", r.URL.String()).
			Int("Status Code", http.StatusUnprocessableEntity).
			Err(err).
			Interface("Account", accountDTO).
			Msg("Failed processing body!")
		return
	}

	_, err = s.AccountService.CreateAccount(accountDTO)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error().
			Str("Method", r.Method).
			Str("Path", r.URL.String()).
			Int("Status Code", http.StatusInternalServerError).
			Err(err).
			Interface("Account", accountDTO).
			Msg("Could not create account!")
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Info().
		Str("Method", r.Method).
		Str("Path", r.URL.String()).
		Int("Status Code", http.StatusCreated).
		Msg("")

}

func (s *AccountServer) Login(w http.ResponseWriter, r *http.Request) {

	log.Info().Str("Method", r.Method).Str("Path", r.URL.String()).Msg("Called endpoint Login!")

	var accountDTO dto.LoginInputDTO

	err := json.NewDecoder(r.Body).Decode(&accountDTO)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)

		log.Info().
			Str("Method", r.Method).
			Str("Path", r.URL.String()).
			Int("Status Code", http.StatusUnprocessableEntity).
			Err(err).
			Interface("Account", accountDTO).
			Msg("Failed processing body!")
		return
	}

	id, err := s.AccountService.Authenticate(accountDTO.CPF, accountDTO.Secret)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		log.Info().
			Str("Method", r.Method).
			Str("Path", r.URL.String()).
			Int("Status Code", http.StatusNotFound).
			Err(err).
			Interface("Account", accountDTO).
			Msg("Account not found!")
		return
	}

	token := s.AuthService.CreateToken(id)

	output := dto.LoginOutputDTO{
		Token: token,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)

	log.Info().
		Str("Method", r.Method).
		Str("Path", r.URL.String()).
		Int("Status Code", http.StatusOK).
		Msg("")

}
