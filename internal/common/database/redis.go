package database

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/yourorg/nms-go/internal/common/config"
)

func NewRedisConnection(cfg config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return rdb, nil
}
