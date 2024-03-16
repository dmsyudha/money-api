package main

import (
	"log"
	"net/http"

	"github.com/dmsyudha/money-api/internal/domain"
	"github.com/dmsyudha/money-api/internal/handler"
	"github.com/dmsyudha/money-api/internal/repository"
	"github.com/dmsyudha/money-api/internal/service"
	"github.com/dmsyudha/money-api/lib/bank"
	db "github.com/dmsyudha/money-api/lib/pq"
	api "github.com/dmsyudha/money-api/pkg/http_client"
	"gorm.io/gorm"
)

type Handler struct {
	AccountHandler  handler.AccountHandler
	TransferHandler handler.TransferHandler
}

func setupDependencies() Handler {
	database, err := db.ConnectToDB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = migrate(database)
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	apiClient := api.NewAPIClient(nil)
	bankAPI := bank.NewBankAPI(apiClient)

	accountRepo := repository.NewAccountRepository(database, bankAPI)
	transferRepo := repository.NewTransferRepository(database, bankAPI)

	accountService := service.NewAccountService(accountRepo)
	transferService := service.NewTransferService(transferRepo, accountRepo)

	accountHandler := handler.NewAccountHandler(accountService)
	transferHandler := handler.NewTransferHandler(transferService)

	return Handler{
		AccountHandler:  accountHandler,
		TransferHandler: transferHandler,
	}
}

func setupRoutes(handler Handler) {
	http.HandleFunc("/api/v1/validate-account", handler.AccountHandler.ValidateAccountHandler())
	http.HandleFunc("/api/v1/create-transfer", handler.TransferHandler.CreateTransferHandler())
}

func migrate(db *gorm.DB) error {
	// AutoMigrate will only add missing fields, won't delete/change existing ones
	return db.AutoMigrate(&domain.Account{}, &domain.Transfer{})
}
