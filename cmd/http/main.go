package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PPAKruNN/golearn/app/handlers"
	"github.com/PPAKruNN/golearn/domain/service"
	"github.com/PPAKruNN/golearn/infra/repository/database"
	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"
)

const (
	PORT = ":5000"
)

func main() {

	// Repository instances
	transferRepo := database.NewTransferRepository()
	fmt.Print("\nTransferrepo\n")

	accountRepo := database.NewAccountRepository()
	fmt.Print("\naccountRepo\n")

	authRepo := database.NewAuthRepository()
	fmt.Print("\nAuthRepo\n")

	// Services instances
	transferService := *service.NewTransferService(transferRepo, accountRepo)
	authService := *service.NewAuthService(authRepo)
	accountService := *service.NewAccountService(accountRepo, authRepo)

	// Handlers instances
	accountServer := handlers.NewAccountServer(transferService, accountService, authService)
	transferServer := handlers.NewTransferServer(transferService, authService)

	// Router
	router := http.NewServeMux()
	router.Handle("/accounts/", accountServer.ServeHTTP())
	router.Handle("/transfers/", transferServer.ServeHTTP())
	router.Handle("/login", http.HandlerFunc(accountServer.Login))

	// Logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger.Info().Msg(fmt.Sprintf("Running server on port %s", PORT))

	log.Fatal(http.ListenAndServe(PORT, router))
}
