package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/PPAKruNN/golearn/app/handlers"
	"github.com/PPAKruNN/golearn/domain/service"
	"github.com/PPAKruNN/golearn/infra/repository/indisk"
	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"
)

const (
	PORT = ":5000"
)

func main() {

	// Placing json files into folder database on current directory.
	CURR_DIR, err := os.Getwd()
	DATABASE_DIR := path.Join(CURR_DIR, "database")

	os.Mkdir(DATABASE_DIR, 0755)

	if err != nil {
		panic("Uai")
	}

	// Repository instances
	transferRepo := indisk.NewTransferRepository(DATABASE_DIR)
	accountRepo := indisk.NewAccountRepository(DATABASE_DIR)
	authRepo := indisk.NewAuthRepository(DATABASE_DIR)

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
