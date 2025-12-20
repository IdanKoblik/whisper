package mock

import (
	"context"
	"strings"
	"testing"
	"whisper-api/config"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TestDBResources struct {
	MongoClient *mongo.Client
	RedisClient *redis.Client
	DB          string
	Collection  string
	Config      *config.Config
}

func Setup(t *testing.T) *TestDBResources {
	cfg := ConfigMock(t)

	ctx := context.Background()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.ConnectionString))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	db := cfg.Mongo.Database
	collection := cfg.Mongo.Collection

	database := mongoClient.Database(db)
	err = database.CreateCollection(ctx, collection)
	if err != nil {
		errStr := err.Error()
		if !strings.Contains(errStr, "already exists") && !strings.Contains(errStr, "NamespaceExists") {
			t.Logf("Collection creation note: %v", err)
		}
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
	})

	return &TestDBResources{
		MongoClient: mongoClient,
		RedisClient: redisClient,
		DB:          db,
		Collection:  collection,
		Config:      cfg,
	}
}

func Teardown(resources *TestDBResources) {
	if resources == nil {
		return
	}

	ctx := context.Background()
	if resources.MongoClient != nil && resources.DB != "" && resources.Collection != "" {
		_ = resources.MongoClient.Database(resources.DB).Collection(resources.Collection).Drop(ctx)
	}
	if resources.RedisClient != nil {
		_ = resources.RedisClient.FlushAll(ctx).Err()
	}
}
