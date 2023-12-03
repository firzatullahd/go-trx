package service

import (
	"context"
	"fmt"
	"go-trx/config"
	aRepository "go-trx/domain/account/repository"
	tError "go-trx/domain/transaction/error"
	"go-trx/domain/transaction/model"
	"go-trx/domain/transaction/repository"
	"go-trx/logger"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type Service interface {
	InsertTransaction(ctx context.Context, payload model.NewTransaction) error
}

type service struct {
	conf        config.Config
	repo        repository.Repository
	accountRepo aRepository.Repository
	redisClient *redis.Client
}

func NewService(conf config.Config, repo repository.Repository, accountRepo aRepository.Repository, redisClient *redis.Client) Service {
	return &service{
		conf:        conf,
		repo:        repo,
		accountRepo: accountRepo,
		redisClient: redisClient,
	}
}

func (s *service) InsertTransaction(ctx context.Context, payload model.NewTransaction) error {

	account, err := s.accountRepo.AccountBalance(ctx, payload.UserID)
	if err != nil {
		logger.Error(ctx, err.Error())
		return err
	}

	key := fmt.Sprintf("lock_%s", payload.ReferenceKey)
	ok, err := s.redisClient.SetNX(ctx, key, "1", time.Duration(s.conf.Constant.TrxTTL)*time.Hour).Result()
	if err != nil {
		logger.Error(ctx, err.Error())
		return err
	}
	if !ok {
		return tError.ErrDuplicateTrx
	}

	logger.Info(ctx, "account %d with balance %v insert transaction %s with amount %v | %v", account.ID, account.Balance, payload.TransactionType, payload.Amount, payload.Amount.Abs())

	if strings.EqualFold(payload.TransactionType, model.TransactionCredit) && payload.Amount.Abs().GreaterThan(account.Balance) {
		return tError.ErrBalanceInsufficient
	}

	tx, err := s.repo.WithTransaction()
	if err != nil {
		logger.Error(ctx, err.Error())
		return err
	}

	err = s.repo.InsertTransaction(ctx, tx, model.AccountTransaction{
		AccountID:       account.ID,
		TransactionType: payload.TransactionType,
		Remark:          payload.Remark,
		Amount:          payload.Amount,
	})
	if err != nil {
		tx.Rollback()
		logger.Error(ctx, err.Error())
		return err
	}

	currentBalance, err := s.repo.CalculateBalance(ctx, tx, account.ID)
	if err != nil {
		tx.Rollback()
		logger.Error(ctx, err.Error())
		return err
	}

	err = s.accountRepo.UpdateBalance(ctx, tx, account.ID, currentBalance)
	if err != nil {
		tx.Rollback()
		logger.Error(ctx, err.Error())
		return err
	}

	tx.Commit()
	return nil
}
