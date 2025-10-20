package db

import (
	"whisper-api/config"

	"github.com/redis/go-redis/v9"
)

func RedisConnection(cfg *config.Config) (*redis.Client) {
    rdb := redis.NewClient(&redis.Options{
        Addr:     cfg.Redis.Addr,
        Password: cfg.Redis.Password,
        DB:       cfg.Redis.DB,  
    })

	return rdb
} 
