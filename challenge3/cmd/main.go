package main

import (
	"log"
	"net/http"

	"challenge3/config"
	"challenge3/handler"
	"challenge3/models"
	"challenge3/repository"
	"challenge3/server"
	"challenge3/service"
)

func main() {
	db, err := config.InitDb()
	db.AutoMigrate(&models.Account{}, &models.Transaction{}, &models.Bank{})
	db.AutoMigrate(&models.Bank{})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database successfully running")
	addr := ":8000"
	mux := http.NewServeMux()

	accountRepo := repository.NewAccountRepo(db)
	accountService := service.NewAccountService(accountRepo)

	bankRepo := repository.NewBankRepo(db)
	bankService := service.NewBankService(bankRepo)

	transactionRepo := repository.NewTransactionRepo(db)
	transactionService := service.NewTransactionService(transactionRepo)

	accountHandler := handler.NewAccountHandler(mux, transactionService, accountService, bankService)
	transactionHandler := handler.NewTransctionHandler(mux, transactionService, accountService)
	bankHandler := handler.NewBankHandler(mux, bankService)

	accountHandler.MapRoutes()
	transactionHandler.MapRoutes()
	bankHandler.MapRoutes()

	defer sqlDB.Close()

	log.Println("Server running on", addr)
	http.ListenAndServe(addr,
		server.ApplicationMiddlewareResponse(
			server.HandleRouteNotFound(mux),
		),
	)
}
