package db

import (
	"context"
	"errors"
	"fmt"
	"time"
	"whisper-api/config"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func RedisConnection(cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	return rdb
}

func DoesExists(cfg *config.Config, key string, fallback string) (bool, error) {
	client := RedisConnection(cfg)
	_, err := client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		collection, err := MongoCollection(cfg)
		if err != nil {
			return false, err
		}

		found, err := FindData(collection, key)
		if err != nil {
			return false, err
		}

		if !found {
			return false, fmt.Errorf("this api token does not exists")
		}

		err = client.Set(ctx, key, fallback, 0).Err()
		if err != nil {
			return false, err
		}

		return true, nil
	}

	if err != nil {
		return false, err
	}

	err = client.Set(ctx, key, fallback, 0).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

func RateLimit(cfg *config.Config, key string, limit int) (int, error) {
	client := RedisConnection(cfg)
	window := time.Minute
	count, err := client.Incr(ctx, key).Result()
	if err != nil {
		return -1, err
	}

	if count == 1 {
		client.Expire(ctx, key, window)
	}

	if count > int64(limit) {
		ttl, _ := client.TTL(ctx, key).Result()
		return int(ttl.Seconds()), nil
	}

	return 0, nil
}
