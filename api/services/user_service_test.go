package services

import (
	"context"
	"testing"
	"os"
	"time"
	"whisper-api/db"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const TEST_DB = "whisper_test"
const TEST_COLLECTION = "users"


func TestRegisterAndUnregisterUser(t *testing.T) {
	os.Setenv("MONGO_CONNECTION", "mongodb://127.0.0.1:27017/?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+2.5.3")
   os.Setenv("WHISPER_DB", "whisper")

   client, err := db.MongoConnection()
   if err != nil {
           t.Fatal(err)
   }

   collection := client.Database(os.Getenv("WHISPER_DB")).Collection("users")
	service := &UserService{Collection: collection}

	user := &User{
		Owner:       "owner1",
		Token:       "token123",
		Subject:     "subject1",
		Subscribers: []string{"sub1", "sub2"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
   defer cancel()
   _, err = collection.DeleteOne(ctx, map[string]string{"_id": user.Owner})
   if err != nil {
      t.Fatalf("failed to cleanup user %s: %v", user.Owner, err)
   }

	// Test RegisterUser success
	err = service.RegisterUser(user)
	assert.NoError(t, err)

	// Test RegisterUser duplicate
	err = service.RegisterUser(user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")

	// Check user exists in DB
	ctx, cancel = context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	count, err := collection.CountDocuments(ctx, bson.M{"_id": user.Owner})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	// Test UnregisterUser success
	err = service.UnregisterUser(user.Owner)
	assert.NoError(t, err)

	// Test UnregisterUser for non-existing user
	err = service.UnregisterUser(user.Owner)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not exists")
}

