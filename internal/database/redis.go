package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sudo-hassan-zahid/go-api-server/internal/config"
	appLogger "github.com/sudo-hassan-zahid/go-api-server/internal/logger"
)

var Rdb *redis.Client

func ConnectRedis(cfg config.RedisConfig) error {
	Rdb = redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:        cfg.Password,
		DB:              cfg.DB,
		DialTimeout:     5 * time.Second,
		ReadTimeout:     3 * time.Second,
		WriteTimeout:    3 * time.Second,
		PoolTimeout:     4 * time.Second,
		MinIdleConns:    2,
		MaxRetries:      3,
		MinRetryBackoff: 100 * time.Millisecond,
		MaxRetryBackoff: 1 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pong, err := Rdb.Ping(ctx).Result()
	if err != nil {
		_ = Rdb.Close()
		Rdb = nil
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	appLogger.Log.Info().Str("pong", pong).Msg("Connected to Redis")
	return nil
}

func CloseRedis() error {
	if Rdb != nil {
		return Rdb.Close()
	}
	return nil
}
