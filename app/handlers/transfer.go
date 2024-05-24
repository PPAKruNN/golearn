package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PPAKruNN/golearn/domain/service"
	"github.com/PPAKruNN/golearn/domain/service/dto"

	"github.com/rs/zerolog/log"
)

type TransferServer struct {
	TransferService service.TransferService
	AuthService     service.AuthService
}

func NewTransferServer(transferService service.TransferService, authService service.AuthService) *TransferServer {
	return &TransferServer{
		TransferService: transferService,
		AuthService:     authService,
	}
}

func (s *TransferServer) ServeHTTP() *http.ServeMux {

	router := http.NewServeMux()
	router.Handle("/transfers/", http.HandlerFunc(s.transferHandler))

	return router

}

// FIX: Change this function name to a more meaningfull name.
func (s *TransferServer) transferHandler(w http.ResponseWriter, r *http.Request) {

	method := r.Method

	fmt.Print(method)

	switch method {
	case http.MethodPost:
		s.CreateTransfer(w, r)
		break
	case http.MethodGet:
		s.ReadTransfers(w, r)
		break

	case "":
		s.ReadTransfers(w, r)
		break

	default:
		w.WriteHeader(http.StatusNotFound)
	}

}

func (s *TransferServer) authorizeAccount(authorization string) (int, error) {

	var token string
	_, scanErr := fmt.Sscanf(authorization, "Bearer %s", &token)
	if scanErr != nil {
		return 0, fmt.Errorf("Invalid bearer token format!")
	}

	accountId, tokenErr := s.AuthService.DecodeToken(token)
	if tokenErr != nil {
		return 0, fmt.Errorf("Invalid token provided!")
	}

	return accountId, nil
}

func (s *TransferServer) ReadTransfers(w http.ResponseWriter, r *http.Request) {

	log.Info().Str("Method", r.Method).Str("Path", r.URL.String()).Msg("Called endpoint ReadTransfers!")

	// Authorization
	authorization := r.Header.Get("Authorization")
	accountId, err := s.authorizeAccount(authorization)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err.Error())

		log.Info().
			Str("Method", r.Method).
			Str("Path", r.URL.String()).
			Int("Status Code", http.StatusUnauthorized).
			Err(err).
			Msg("Failed authorizing request!")
		return
	}

	transfers := s.TransferService.ReadTransfersByAccount(accountId)

	json.NewEncoder(w).Encode(transfers)
	log.Info().
		Str("Method", r.Method).
		Str("Path", r.URL.String()).
		Int("Status Code", http.StatusOK).
		Msg("")

}

func (s *TransferServer) CreateTransfer(w http.ResponseWriter, r *http.Request) {

	log.Info().Str("Method", r.Method).Str("Path", r.URL.String()).Msg("Called endpoint CreateTransfer!")

	// Authorization
	authorization := r.Header.Get("Authorization")
	accountId, authErr := s.authorizeAccount(authorization)
	if authErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(authErr.Error())

		log.Info().
			Str("Method", r.Method).
			Str("Path", r.URL.String()).
			Int("Status Code", http.StatusUnauthorized).
			Err(authErr).
			Msg("Failed authorizing request!")
		return
	}

	var input dto.CreateTrasnferInputDTO
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(err.Error())

		log.Info().
			Str("Method", r.Method).
			Str("Path", r.URL.String()).
			Int("Status Code", http.StatusUnprocessableEntity).
			Err(err).
			Msg("Failed processing body!")
		return
	}

	// FIX: DTO should not ask for OriginID on json.
	// Temporary solution: Force origin to be accountId.
	if input.AccountOriginID != accountId {
		log.Warn().
			Str("Method", r.Method).
			Str("Path", r.URL.String()).
			Interface("UserSentInput", input).
			Int("NewCorrectedOriginID", accountId).
			Msg("Manually correcting an transfer originID to a accountID obtained using account authentication!")

		input.AccountOriginID = accountId
	}

	statusCode, transferErr := s.TransferService.CreateTransfer(input)
	if transferErr != nil {
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(transferErr.Error())
		log.Info().
			Str("Method", r.Method).
			Str("Path", r.URL.String()).
			Int("Status Code", statusCode).
			Err(transferErr).
			Msg("Failed creating a transfer!")
		return
	}

	w.WriteHeader(http.StatusCreated)

	log.Info().
		Str("Method", r.Method).
		Str("Path", r.URL.String()).
		Interface("Account", input).
		Int("Status Code", http.StatusCreated).
		Msg("")
}
