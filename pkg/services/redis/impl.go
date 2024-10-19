package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/rafaalrazzak/e-campus-be/pkg/framework/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ECampusRedisDB struct {
	Client *redis.Client
}

type RedisParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    config.Config
	Logger    *zap.Logger
}

func NewECampusRedisDBImpl(client *redis.Client) *ECampusRedisDB {
	return &ECampusRedisDB{Client: client}
}

func NewRedisConn(p RedisParams) (*redis.Client, error) {
	opts, err := redis.ParseURL(p.Config.Redis)
	if err != nil {
		p.Logger.Error("Failed to parse Redis URL", zap.Error(err))
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return connectRedis(ctx, client, p.Logger)
		},
		OnStop: func(ctx context.Context) error {
			p.Logger.Info("Closing Redis connection")
			return client.Close()
		},
	})

	return client, nil
}

func connectRedis(ctx context.Context, client *redis.Client, logger *zap.Logger) error {
	logger.Info("Attempting to connect to Redis")

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := client.Ping(ctxWithTimeout).Err()
	if err != nil {
		logger.Error("Failed to connect to Redis", zap.Error(err))
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Successfully connected to Redis")
	return nil
}

func InitializeRedis(redisDB *ECampusRedisDB, logger *zap.Logger) {
	logger.Info("Initializing Redis connection")
	err := redisDB.Client.Ping(context.Background()).Err()
	if err != nil {
		logger.Error("Failed to ping Redis", zap.Error(err))
	} else {
		logger.Info("Redis connection initialized successfully")
	}
}
