package service

import (
	"context"
	"errors"
	"go-trx/config"
	"go-trx/domain/account/model"
	"go-trx/domain/account/repository"
	"go-trx/logger"
)

type Service interface {
	AccountBalance(ctx context.Context, userID uint64) (*model.Account, error)
}

type service struct {
	conf config.Config
	repo repository.Repository
}

func NewService(conf config.Config, repo repository.Repository) Service {
	return &service{
		conf: conf,
		repo: repo,
	}
}

func (s *service) AccountBalance(ctx context.Context, userID uint64) (*model.Account, error) {
	account, err := s.repo.AccountBalance(ctx, userID)
	if err != nil {
		logger.Error(ctx, err.Error())
		return nil, err
	}
	if account == nil {
		return nil, errors.New("account not found")
	}
	return account, nil
}
