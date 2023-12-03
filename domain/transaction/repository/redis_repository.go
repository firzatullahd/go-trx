package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

//go:generate mockgen -source=redis_repository.go -destination=./mock/redis_repository.go -package=repository
type RedisRepository interface {
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
}

type redisRepository struct {
	redis *redis.Client
}

func NewRedisRepository(redisClient *redis.Client) RedisRepository {
	return &redisRepository{
		redis: redisClient,
	}
}

func (r *redisRepository) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.redis.SetNX(ctx, key, value, expiration).Result()
}
