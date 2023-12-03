package server

import (
	"context"
	"fmt"
	"go-trx/config"
	"go-trx/db"
	a "go-trx/domain/account"
	aRepository "go-trx/domain/account/repository"
	aService "go-trx/domain/account/service"
	hc "go-trx/domain/health_check"
	hcRepository "go-trx/domain/health_check/repository"
	hcService "go-trx/domain/health_check/service"
	t "go-trx/domain/transaction"
	tRepository "go-trx/domain/transaction/repository"
	tService "go-trx/domain/transaction/service"
	"go-trx/logger"
	"log"
	"os"
	"os/signal"

	"github.com/labstack/echo/v4"
)

func Start(conf config.Config) {
	ctx := context.Background()

	masterPSQL, slavePSQL := db.InitializePSQL(conf.Postgres)
	redisClient := db.InitializeRedisClient(conf.Redis)

	hcRepo := hcRepository.NewRepository(masterPSQL, slavePSQL)
	hcSvc := hcService.NewService(conf, hcRepo)
	hcHandler := hc.NewHandler(hcSvc)

	aRepo := aRepository.NewRepository(masterPSQL, slavePSQL)
	aSvc := aService.NewService(conf, aRepo)
	aHandler := a.NewHandler(aSvc)

	tRepo := tRepository.NewRepository(masterPSQL, slavePSQL)
	tRedisRepo := tRepository.NewRedisRepository(redisClient)
	tSvc := tService.NewService(conf, tRepo, aRepo, tRedisRepo)
	tHandler := t.NewHandler(tSvc)

	e := echo.New()
	e = InitializeRouter(e, hcHandler, aHandler, tHandler)
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		logger.Info(ctx, "We received an interrupt signal, shut down.")
		if err := e.Shutdown(context.Background()); err != nil {
			logger.Error(ctx, "HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
		logger.Info(ctx, "Bye.")
	}()

	log.Fatal(e.Start(fmt.Sprintf(":%s", conf.App.Port)))
}
