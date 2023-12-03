package service

import (
	"context"
	"errors"
	"fmt"
	"go-trx/config"
	aRepository "go-trx/domain/account/repository"
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

	key := fmt.Sprintf("lock_%s", payload.ReferenceKey)

	ok, err := s.redisClient.SetNX(ctx, key, "1", time.Duration(s.conf.Constant.TrxTTL)*time.Hour).Result()
	if err != nil {
		logger.Error(ctx, err.Error())
		return err
	}
	if !ok {
		return errors.New("duplicate transaction")
	}

	account, err := s.accountRepo.AccountBalance(ctx, payload.UserID)
	if err != nil {
		logger.Error(ctx, err.Error())
		return err
	}

	if strings.EqualFold(payload.TransactionType, model.TransactionCredit) && !account.Balance.GreaterThanOrEqual(payload.Amount) {
		return errors.New("balance insufficient")
	}

	err = s.repo.InsertTransaction(ctx, model.AccountTransaction{
		AccountID:       account.ID,
		TransactionType: payload.TransactionType,
		Remark:          payload.Remark,
		Amount:          payload.Amount,
	})
	if err != nil {
		logger.Error(ctx, err.Error())
		return err
	}

	currentBalance, err := s.repo.CalculateBalance(ctx, account.ID)
	if err != nil {
		logger.Error(ctx, err.Error())
		return err
	}
	err = s.accountRepo.UpdateBalance(ctx, account.ID, currentBalance)
	if err != nil {
		logger.Error(ctx, err.Error())
		return err
	}

	return nil
}
