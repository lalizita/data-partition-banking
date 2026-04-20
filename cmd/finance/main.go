package main

import (
	"context"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/lalizita/shard-banking/internal/config"
	"github.com/lalizita/shard-banking/internal/infraestructure/db"
	txHandlers "github.com/lalizita/shard-banking/internal/services/transaction/handlers"
	txRepository "github.com/lalizita/shard-banking/internal/services/transaction/repository"
	txService "github.com/lalizita/shard-banking/internal/services/transaction/service"
)

func main() {
	ctx := context.Background()
	e := echo.New()

	cfg, err := config.Load(ctx)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	financeDB, err := db.NewFinanceDB(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to connect finance database: %v", err)
	}

	txRepo := txRepository.NewTransactionRepository(financeDB)
	txSvc := txService.NewTransactionService(txRepo)
	txHandler := txHandlers.NewTransactionHandler(e, txSvc)
	txHandler.RegisterRoutes()

	port := os.Getenv("PORT_FINANCE")
	if port == "" {
		port = "8081"
	}

	addr := ":" + port
	log.Printf("starting finance server on %s", addr)
	if err := e.Start(addr); err != nil {
		log.Fatal(err)
	}
}
