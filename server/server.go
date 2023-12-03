package server

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"go-trx/config"
	"go-trx/db"
	hc "go-trx/domain/health_check"
	hcRepository "go-trx/domain/health_check/repository"
	hcService "go-trx/domain/health_check/service"
	"go-trx/logger"
	"log"
	"os"
	"os/signal"
)

func Start(conf config.Config) {
	ctx := context.Background()

	masterPSQL, slavePSQL := db.InitializePSQL(conf.Postgres)
	redisClient := db.InitializeRedisClient(conf.Redis)

	fmt.Println(masterPSQL, slavePSQL, redisClient)

	hcRepo := hcRepository.NewRepository(masterPSQL, slavePSQL)
	hcSvc := hcService.NewService(conf, hcRepo)
	hcHandler := hc.NewHealthCheckHandler(hcSvc)

	e := echo.New()
	//todo: middleware
	//todo: setup router
	e = InitializeRouter(e, hcHandler)
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
