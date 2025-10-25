package db

import (
	"context"
	"testing"
	"whisper-api/mock"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestMongoConnection(t *testing.T) {
	cfg := mock.ConfigMock(t)

	client, err := MongoConnection(cfg)
	if err != nil {
		t.Fatalf("failed to connect to mongo: %v", err)
	}
	if client == nil {
		t.Fatal("expected mongo client, got nil")
	}
}

func TestMongoCollection(t *testing.T) {
	cfg := mock.ConfigMock(t)

	collection, err := MongoCollection(cfg)
	if err != nil {
		t.Fatalf("failed to get collection: %v", err)
	}
	if collection == nil {
		t.Fatal("expected collection, got nil")
	}
}

func TestInsertAndFindData(t *testing.T) {
	cfg := mock.ConfigMock(t)
	collection, err := MongoCollection(cfg)
	if err != nil {
		t.Fatalf("failed to get collection: %v", err)
	}

	testKey := "test_user_1"
	testData := bson.M{"_id": testKey, "name": "Test User"}

	collection.DeleteOne(context.Background(), bson.M{"_id": testKey})
	defer collection.DeleteOne(context.Background(), bson.M{"_id": testKey})

	// Test InsertData
	if err := InsertData(cfg, testData, testKey); err != nil {
		t.Fatalf("InsertData failed: %v", err)
	}

	found, err := FindData(collection, testKey)
	if err != nil {
		t.Fatalf("FindData failed: %v", err)
	}
	if !found {
		t.Fatal("expected document to exist after insertion")
	}

	if err := InsertData(cfg, testData, testKey); err == nil {
		t.Fatal("expected duplicate insertion to fail")
	}
}

func TestDeleteData(t *testing.T) {
	cfg := mock.ConfigMock(t)
	collection, err := MongoCollection(cfg)
	if err != nil {
		t.Fatalf("failed to get collection: %v", err)
	}

	testKey := "test_user_2"
	testData := bson.M{"_id": testKey, "name": "Delete User"}

	if err := InsertData(cfg, testData, testKey); err != nil {
		t.Fatalf("InsertData failed: %v", err)
	}

	if err := DeleteData(cfg, testKey); err != nil {
		t.Fatalf("DeleteData failed: %v", err)
	}

	found, err := FindData(collection, testKey)
	if err != nil {
		t.Fatalf("FindData failed: %v", err)
	}
	if found {
		t.Fatal("expected document to be deleted")
	}

	if err := DeleteData(cfg, testKey); err == nil {
		t.Fatal("expected error when deleting non-existent document")
	}
}
