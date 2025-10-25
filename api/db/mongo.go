package db

import (
	"context"
	"fmt"
	"time"
	"whisper-api/config"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const TIMEOUT = 5 * time.Second

func MongoConnection(cfg *config.Config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(cfg.Mongo.ConnectionString))
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, err
}

func MongoCollection(cfg *config.Config) (*mongo.Collection, error) {
	client, err := MongoConnection(cfg)
	if err != nil {
		return nil, err
	}

	return client.Database(cfg.Mongo.Database).Collection("users"), nil
}

func FindData(collection *mongo.Collection, key string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	count, err := collection.CountDocuments(ctx, bson.M{"_id": key})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func InsertData(cfg *config.Config, data interface{}, key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	collection, err := MongoCollection(cfg)
	if err != nil {
		return err
	}

	found, err := FindData(collection, key)
	if err != nil {
		return err
	}

	if found {
		return fmt.Errorf("dupelicated data")
	}

	_, err = collection.InsertOne(ctx, data)
	return err
}

func DeleteData(cfg *config.Config, key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	collection, err := MongoCollection(cfg)
	if err != nil {
		return err
	}

	found, err := FindData(collection, key)
	if err != nil {
		return err
	}

	if !found {
		return fmt.Errorf("this document does not exists")
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": key})
	return err
}
