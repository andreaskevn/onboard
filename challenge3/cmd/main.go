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
	// "go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

func main() {
	err := config.InitLog()
	if err != nil {
		panic(err)
	}

	defer config.Log.Sync()
	config.Log.Info("Logger initialized")

	db, err := config.InitDb()
	db.AutoMigrate(&models.Account{}, &models.Transaction{}, &models.Bank{})
	db.AutoMigrate(&models.Bank{})
	db.AutoMigrate(&models.IdempotencyKey{})
	if err != nil {
		config.Log.Fatal("Failed to connect database", zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		config.Log.Fatal("Failed to get SQL DB", zap.Error(err))
	}

	config.Log.Info("Database Initialized")

	redis, err := config.InitRedis()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Redis Initialized")
	defer redis.Close()

	cleanup := config.InitOpenTelemetry("bank-service")
	defer cleanup()
	
	server.InitMetrics()

	addr := ":8000"
	mux := http.NewServeMux()

	accountRepo := repository.NewAccountRepo(db)
	bankRepo := repository.NewBankRepo(db)
	transactionRepo := repository.NewTransactionRepo(db)

	bankService := service.NewBankService(bankRepo)
	accountService := service.NewAccountService(accountRepo, bankRepo)
	transactionService := service.NewTransactionService(transactionRepo, accountRepo)

	accountHandler := handler.NewAccountHandler(mux, transactionService, accountService, bankService, redis)
	transactionHandler := handler.NewTransctionHandler(mux, transactionService, accountService, redis)
	bankHandler := handler.NewBankHandler(mux, bankService)

	accountHandler.MapRoutes()
	transactionHandler.MapRoutes()
	bankHandler.MapRoutes()

	defer sqlDB.Close()

	handlerChain := server.ApplicationMiddlewareResponse(
		server.HandleRouteNotFound(mux),
	)
	handlerChain = server.IdempotencyMiddlewareRedis(redis, handlerChain)
	handlerChain = server.MetricsMiddleware(handlerChain)
	config.Log.Info("Server starting", zap.String("address", addr))

	err = http.ListenAndServe(addr, handlerChain)
}
