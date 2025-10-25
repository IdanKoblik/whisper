package db

import (
	"context"
	"testing"
	"time"
	"whisper-api/mock"

	"go.mongodb.org/mongo-driver/bson"
)

func TestDoesExists(t *testing.T) {
	cfg := mock.ConfigMock(t)
	client := RedisConnection(cfg)
	defer client.FlushDB(context.Background())

	collection, err := MongoCollection(cfg)
	if err != nil {
		t.Fatalf("failed to get mongo collection: %v", err)
	}

	testKey := "test_token"
	testData := bson.M{"_id": testKey, "name": "Test User"}
	collection.InsertOne(context.Background(), testData)
	defer collection.DeleteOne(context.Background(), bson.M{"_id": testKey})

	exists, err := DoesExists(cfg, testKey, "cached_value")
	if err != nil {
		t.Fatalf("DoesExists failed: %v", err)
	}
	if !exists {
		t.Fatal("expected key to exist after checking Mongo fallback")
	}

	val, _ := client.Get(context.Background(), testKey).Result()
	if val != "cached_value" {
		t.Fatalf("expected Redis key to be set with fallback value, got %s", val)
	}

	exists, err = DoesExists(cfg, testKey, "new_value")
	if err != nil {
		t.Fatalf("DoesExists failed: %v", err)
	}
	if !exists {
		t.Fatal("expected key to exist in Redis")
	}

	val, _ = client.Get(context.Background(), testKey).Result()
	if val != "new_value" {
		t.Fatalf("expected Redis key to be updated to new value, got %s", val)
	}

	notExistKey := "missing_token"
	exists, err = DoesExists(cfg, notExistKey, "some_value")
	if err == nil {
		t.Fatal("expected error for missing key in both Redis and Mongo")
	}
	if exists {
		t.Fatal("expected exists to be false for missing key")
	}
}

func TestRateLimit(t *testing.T) {
	cfg := mock.ConfigMock(t)
	client := RedisConnection(cfg)
	defer client.FlushDB(context.Background())

	key := "rate_limit_test"
	limit := 3

	for i := 1; i <= limit; i++ {
		ttl, err := RateLimit(cfg, key, limit)
		if err != nil {
			t.Fatalf("RateLimit failed: %v", err)
		}
		if ttl != 0 {
			t.Fatalf("expected ttl 0 before limit exceeded, got %d", ttl)
		}
	}

	ttl, err := RateLimit(cfg, key, limit)
	if err != nil {
		t.Fatalf("RateLimit failed: %v", err)
	}
	if ttl <= 0 {
		t.Fatalf("expected positive TTL after limit exceeded, got %d", ttl)
	}

	time.Sleep(time.Second * 61)
	count, _ := client.Get(context.Background(), key).Result()
	if count != "" {
		t.Fatalf("expected Redis key to expire, got %s", count)
	}
}
