package db

import (
	"context"
	"testing"
	"time"
	"whisper-api/mock"
)

const RATE_LIMIT_DURATION = time.Second * 10

func TestRateLimit(t *testing.T) {
	resources := mock.Setup(t)
	defer mock.Teardown(resources)

	client := resources.RedisClient

	key := "rate_limit_test"
	limit := 3

	for i := 1; i <= limit; i++ {
		ttl, err := RateLimit(client, key, limit, RATE_LIMIT_DURATION)
		if err != nil {
			t.Fatalf("RateLimit failed: %v", err)
		}
		if ttl != 0 {
			t.Fatalf("expected ttl 0 before limit exceeded, got %d", ttl)
		}
	}

	ttl, err := RateLimit(client, key, limit, RATE_LIMIT_DURATION)
	if err != nil {
		t.Fatalf("RateLimit failed: %v", err)
	}
	if ttl <= 0 {
		t.Fatalf("expected positive TTL after limit exceeded, got %d", ttl)
	}

	time.Sleep(RATE_LIMIT_DURATION + time.Second)
	count, _ := client.Get(context.Background(), key).Result()
	if count != "" {
		t.Fatalf("expected Redis key to expire, got %s", count)
	}
}
