package db

import (
	"github.com/go-redis/redis/v8"
	"go-trx/config"
)

func InitializeRedisClient(conf config.Redis) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     conf.Host,
		Password: conf.Password,
		DB:       conf.DB,
	})
}
