package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/restlesswhy/eth-balance-searcher/config"
)

func New(cfg *config.Config, ctx context.Context) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "api-redis:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, errors.Wrap(err, "check connection to redis error")
	}

	return client, nil
}
