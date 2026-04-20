package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lalizita/shard-banking/internal/config"
	"github.com/lalizita/shard-banking/internal/infraestructure/db"
	accountHandlers "github.com/lalizita/shard-banking/internal/services/account/handlers"
	"github.com/lalizita/shard-banking/internal/services/account/repository"
	accService "github.com/lalizita/shard-banking/internal/services/account/service"
	"github.com/lalizita/shard-banking/pkg/sharding"
)

func main() {
	ctx := context.Background()
	e := echo.New()
	cfg, err := config.Load(ctx)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	shardUnits, err := strconv.Atoi(cfg.TransactionShards)
	if err != nil {
		log.Fatal("failed to transform shard units")
	}

	shardRouter := sharding.NewShardingRouter(shardUnits)

	accountDB, err := db.NewAccountDB(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to connect account database %v", err)
	}

	accRepo := repository.NewAccountRepository(accountDB, shardRouter)
	accSvc := accService.NewAccountService(accRepo)
	accHandler := accountHandlers.NewAccountHandler(e, accSvc)
	accHandler.RegisterRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("starting web server on %s", addr)
	if err := e.Start(addr); err != nil {
		log.Fatal(err)
	}
}
