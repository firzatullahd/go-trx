package db

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go-trx/config"
	"go-trx/logger"
	"time"
)

func InitializePSQL(conf config.Postgres) (master *sqlx.DB, slave *sqlx.DB) {
	var err error
	ctx := context.Background()
	master, err = sqlx.Open("postgres", conf.Master.ConnectionString())
	if err != nil {
		logger.Fatal(ctx, "Can't connect to master DB %+v", err)
	}

	master.SetMaxIdleConns(conf.MaxIdleCons)
	master.SetMaxOpenConns(conf.MaxOpenCons)
	master.SetConnMaxLifetime(time.Duration(conf.ConMaxLifetime) * time.Millisecond)
	master.SetConnMaxIdleTime(time.Duration(conf.ConMaxIdleTime) * time.Millisecond)

	slave, err = sqlx.Connect("postgres", conf.Slave.ConnectionString())
	if err != nil {
		logger.Fatal(ctx, "Can't connect to slave DB %+v", err)
	}

	slave.SetMaxIdleConns(conf.MaxIdleCons)
	slave.SetMaxOpenConns(conf.MaxOpenCons)
	slave.SetConnMaxLifetime(time.Duration(conf.ConMaxLifetime) * time.Millisecond)
	slave.SetConnMaxIdleTime(time.Duration(conf.ConMaxIdleTime) * time.Millisecond)

	return
}
