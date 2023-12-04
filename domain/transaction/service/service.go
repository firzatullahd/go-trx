package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-trx/config"
	aRepository "go-trx/domain/account/repository"
	tError "go-trx/domain/transaction/error"
	"go-trx/domain/transaction/model"
	"go-trx/domain/transaction/repository"
	"go-trx/logger"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type Service interface {
	InsertTransaction(ctx context.Context, payload model.NewTransaction) error
}

type service struct {
	conf        config.Config
	repo        repository.Repository
	accountRepo aRepository.Repository
	redisRepo   repository.RedisRepository
}

func NewService(conf config.Config, repo repository.Repository, accountRepo aRepository.Repository, redisRepo repository.RedisRepository) Service {
	return &service{
		conf:        conf,
		repo:        repo,
		accountRepo: accountRepo,
		redisRepo:   redisRepo,
	}
}

func (s *service) InsertTransaction(ctx context.Context, payload model.NewTransaction) error {

	account, errAccount := s.accountRepo.AccountBalance(ctx, payload.UserID)
	if errAccount != nil && !errors.Is(errAccount, sql.ErrNoRows) {
		logger.Error(ctx, errAccount.Error())
		return errAccount
	}

	key := fmt.Sprintf("lock_%s", payload.ReferenceKey)
	ok, err := s.redisRepo.SetNX(ctx, key, "1", time.Duration(s.conf.Constant.TrxTTL)*time.Hour)
	if err != nil {
		logger.Error(ctx, err.Error())
		return err
	}
	if !ok {
		return tError.ErrDuplicateTrx
	}

	logger.Info(ctx, "accountID %d with balance %v insert transaction %s with amount %v", account.ID, account.Balance, payload.TransactionType, payload.Amount)

	if strings.EqualFold(payload.TransactionType, model.TransactionCredit) && payload.Amount.Abs().GreaterThan(account.Balance) {
		return tError.ErrBalanceInsufficient
	}

	if err := s.repo.WithTransaction(ctx, func(tx *sqlx.Tx) error {

		if errors.Is(errAccount, sql.ErrNoRows) {
			logger.Info(ctx, "insert new account with userID %d", payload.UserID)
			account, err = s.accountRepo.InsertAccount(ctx, tx, payload.UserID)
			if err != nil {
				logger.Error(ctx, err.Error())
				return err
			}
		}

		if err := s.repo.InsertTransaction(ctx, tx, model.AccountTransaction{
			AccountID:       account.ID,
			TransactionType: payload.TransactionType,
			Remark:          payload.Remark,
			Amount:          payload.Amount,
		}); err != nil {
			logger.Error(ctx, err.Error())
			return err
		}

		currentBalance, err := s.repo.CalculateBalance(ctx, tx, account.ID)
		if err != nil {
			tx.Rollback()
			logger.Error(ctx, err.Error())
			return err
		}

		if err := s.accountRepo.UpdateBalance(ctx, tx, account.ID, currentBalance); err != nil {
			logger.Error(ctx, err.Error())
			return err
		}

		return nil

	}); err != nil {
		logger.Error(ctx, err.Error())
		return err
	}

	return nil
}
