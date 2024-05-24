package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PPAKruNN/golearn/domain/service"
	"github.com/PPAKruNN/golearn/domain/service/dto"
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

	fmt.Printf("Chamou read Transfer! \n")

	// Authorization
	authorization := r.Header.Get("Authorization")
	accountId, err := s.authorizeAccount(authorization)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	transfers := s.TransferService.ReadTransfersByAccount(accountId)

	json.NewEncoder(w).Encode(transfers)

}

func (s *TransferServer) CreateTransfer(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("Chamou create Transfer! \n")

	// Authorization
	authorization := r.Header.Get("Authorization")
	accountId, authErr := s.authorizeAccount(authorization)
	if authErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(authErr.Error())
		return
	}

	var input dto.CreateTrasnferInputDTO
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// FIX: DTO should not ask for OriginID on json.
	// Temporary solution: Force origin to be accountId.
	input.AccountOriginID = accountId

	statusCode, transferErr := s.TransferService.CreateTransfer(input)
	if transferErr != nil {
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(transferErr.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)

}
