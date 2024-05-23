package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PPAKruNN/golearn/app/handlers"
	"github.com/PPAKruNN/golearn/domain/service"
	"github.com/PPAKruNN/golearn/infra/repository/inmemory"
)

func main() {
	// Repository instances
	transferRepo := inmemory.NewTransferRepository()
	accountRepo := inmemory.NewAccountRepository()
	authRepo := inmemory.NewAuthRepository()

	// Services instances
	transferService := *service.NewTransferService(transferRepo, accountRepo)
	authService := *service.NewAuthService(authRepo)
	accountService := *service.NewAccountService(accountRepo, authRepo)

	// Handlers instances
	accountServer := handlers.NewAccountServer(transferService, accountService, authService)
	transferServer := handlers.NewTransferServer(transferService, authService)

	router := http.NewServeMux()
	router.Handle("/", accountServer.ServeHTTP())
	router.Handle("/", transferServer.ServeHTTP())

	fmt.Println("Running Server!")
	log.Fatal(http.ListenAndServe(":5000", router))

}
