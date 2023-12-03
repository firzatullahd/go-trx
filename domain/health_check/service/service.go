package service

import (
	"context"
	"go-trx/config"
	"go-trx/domain/health_check/repository"
	"go-trx/logger"
)

type Service interface {
	Ping(ctx context.Context) error
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

func (s *service) Ping(ctx context.Context) error {
	err := s.repo.Ping(ctx)
	if err != nil {
		logger.Error(ctx, err.Error())
		return err
	}
	return nil
}
