package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"go-trx/logger"
)

type Repository interface {
	Ping(ctx context.Context) error
}

type repository struct {
	masterPSQL *sqlx.DB
	slavePSQL  *sqlx.DB
}

func NewRepository(masterPSQL *sqlx.DB, slavePSQL *sqlx.DB) Repository {
	return &repository{
		masterPSQL: masterPSQL,
		slavePSQL:  slavePSQL,
	}
}

func (r *repository) Ping(ctx context.Context) error {
	err := r.slavePSQL.Ping()
	if err != nil {
		logger.Error(ctx, err.Error())
		return err
	}
	err = r.masterPSQL.Ping()
	if err != nil {
		logger.Error(ctx, err.Error())
		return err
	}
	return nil
}
