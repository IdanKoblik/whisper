package db

import (
	"time"
	"context"
	"whisper-api/config"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const TIMEOUT = 5 * time.Second

func MongoConnection(cfg *config.Config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(cfg.Mongo.ConnectionURL))
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, err
}

