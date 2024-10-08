package redis

import (
	"context"
	"github.com/rafaalrazzak/e-campus-be/pkg/framework/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

type ECampusRedis struct {
	Client *redis.Client
}

func NewRedisImpl(lc fx.Lifecycle, cfg config.Config) *ECampusRedis {
	opts, err := redis.ParseURL(cfg.Redis)
	if err != nil {
		panic("Failed to parse Redis URL: " + err.Error())
	}

	client := redis.NewClient(opts)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := client.Ping(ctx).Err(); err != nil {
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return &ECampusRedis{Client: client}
}
