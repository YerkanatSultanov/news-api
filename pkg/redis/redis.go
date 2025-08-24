package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"news-api/internal/config"
	"news-api/pkg/logger"
)

func InitRedis(cfg config.RedisConfig) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       0,
	})

	ctx := context.Background()
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		logger.Log.Error("Failed to connect to Redis", "address", cfg.Addr, "error", err)
		return nil
	}

	logger.Log.Info("Connected to Redis successfully", "address", cfg.Addr)
	return redisClient
}
