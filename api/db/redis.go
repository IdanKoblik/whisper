package db

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

const RATE_LIMIT_PREFIX = "limit:"

func RateLimit(client *redis.Client, key string, limit int, window time.Duration) (int, error) {
	count, err := client.Incr(ctx, RATE_LIMIT_PREFIX+key).Result()
	if err != nil {
		return -1, err
	}

	if count == 1 {
		client.Expire(ctx, RATE_LIMIT_PREFIX+key, window)
	}

	if count > int64(limit) {
		ttl, _ := client.TTL(ctx, RATE_LIMIT_PREFIX+key).Result()
		return int(ttl.Seconds()), nil
	}

	return 0, nil
}
