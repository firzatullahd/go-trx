package db

import (
	"context"
	"fmt"
	"go-trx/config"
	"go-trx/logger"

	"github.com/go-redis/redis/v8"
)

func InitializeRedisClient(conf config.Redis) *redis.Client {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Fatal(ctx, "Can' Ping Redis Client %v", err)
	}

	return client
}
