package repository

import (
	"context"
	"go-trx/domain/account/model"
	"go-trx/logger"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -source=repository.go -destination=./mock/repository.go -package=repository
type Repository interface {
	AccountBalance(ctx context.Context, userID uint64) (model.Account, error)
	UpdateBalance(ctx context.Context, tx *sqlx.Tx, accountID uint64, balance float64) error
	InsertAccount(ctx context.Context, tx *sqlx.Tx, userID uint64) (model.Account, error)
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

func (r *repository) AccountBalance(ctx context.Context, userID uint64) (model.Account, error) {
	var result model.Account

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	query, args, err := sq.Select(`id, user_id, balance, created_at, updated_at, deleted_at`).From(`account`).Where(`deleted_at isnull`).Where(squirrel.Eq{`user_id`: userID}).ToSql()
	if err != nil {
		logger.Error(ctx, err.Error())
		return result, err
	}
	err = r.slavePSQL.QueryRowxContext(ctx, query, args...).StructScan(&result)
	if err != nil {
		logger.Error(ctx, err.Error())
		return result, err
	}

	return result, nil
}

func (r *repository) UpdateBalance(ctx context.Context, tx *sqlx.Tx, accountID uint64, balance float64) error {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	query, args, err := sq.Update(`account`).
		Set(`updated_at`, time.Now()).
		Set(`balance`, balance).
		Where(`deleted_at isnull`).
		Where(squirrel.Eq{`id`: accountID}).ToSql()
	if err != nil {
		logger.Error(ctx, err.Error())
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		logger.Error(ctx, err.Error())
		return err
	}

	return nil
}

func (r *repository) InsertAccount(ctx context.Context, tx *sqlx.Tx, userID uint64) (model.Account, error) {
	var result model.Account

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	query, args, err := sq.Insert(`account`).Columns(`user_id`, `balance`).Values(userID, 0).Suffix(`returning id, user_id, balance, created_at, updated_at, deleted_at`).ToSql()
	if err != nil {
		logger.Error(ctx, err.Error())
		return result, err
	}

	err = tx.QueryRowxContext(ctx, query, args...).StructScan(&result)
	if err != nil {
		logger.Error(ctx, err.Error())
		return result, err
	}

	return result, nil

}
