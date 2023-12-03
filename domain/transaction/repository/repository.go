package repository

import (
	"context"
	"database/sql"
	"errors"
	"go-trx/domain/transaction/model"
	"go-trx/logger"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	WithTransaction() (*sqlx.Tx, error)
	InsertTransaction(ctx context.Context, tx *sqlx.Tx, payload model.AccountTransaction) error
	CalculateBalance(ctx context.Context, tx *sqlx.Tx, accountID uint64) (float64, error)
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

func (r *repository) InsertTransaction(ctx context.Context, tx *sqlx.Tx, payload model.AccountTransaction) error {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	query, args, err := sq.Insert(`account_transaction`).Columns(`account_id`, `transaction_type`, `remark`, `amount`).Values(payload.AccountID, payload.TransactionType, payload.Remark, payload.Amount).ToSql()
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

func (r *repository) CalculateBalance(ctx context.Context, tx *sqlx.Tx, accountID uint64) (float64, error) {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	query, args, err := sq.Select(`COALESCE(SUM(amount), 0)`).From(`account_transaction`).Where(squirrel.Eq{`account_id`: accountID}).ToSql()
	if err != nil {
		logger.Error(ctx, err.Error())
		return 0, err
	}
	var balance float64
	err = tx.QueryRowxContext(ctx, query, args...).Scan(&balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		logger.Error(ctx, err.Error())
		return 0, err
	}

	return balance, nil
}

func (r *repository) WithTransaction() (*sqlx.Tx, error) {
	tx, err := r.masterPSQL.Beginx()
	if err != nil {
		return nil, err
	}

	return tx, nil
}
