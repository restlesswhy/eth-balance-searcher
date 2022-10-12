package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type cache struct {
	client *redis.Client
}

func New(client *redis.Client) *cache {
	return &cache{client: client}
}

func (c *cache) Get(ctx context.Context, key string, dest interface{}) (bool, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, errors.Wrap(err, "get from cache error")
	}

	return true, json.Unmarshal([]byte(val), dest)
}

func (c *cache) Set(ctx context.Context, key string, value interface{}) error {
	p, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "marshal data error")
	}

	return c.client.Set(ctx, key, p, time.Hour*1).Err()
}
