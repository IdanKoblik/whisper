package services

import (
	"context"
	"os"
	"testing"
	"time"
	"whisper-api/db"
	"whisper-api/mock"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const TEST_DB = "whisper_test"
const TEST_COLLECTION = "users"


func TestRegisterAndUnregisterUser(t *testing.T) {
	cfg := mock.ConfigMock(t)		

	client, err := db.MongoConnection(&cfg)
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
   _, err = collection.DeleteOne(ctx, map[string]string{"_id": user.Token})
   if err != nil {
      t.Fatalf("failed to cleanup user %s: %v", user.Owner, err)
   }

	err = service.RegisterUser(user)
	assert.NoError(t, err)

	err = service.RegisterUser(user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")

	ctx, cancel = context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	count, err := collection.CountDocuments(ctx, bson.M{"_id": user.Token})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	err = service.UnregisterUser(user.Token)
	assert.NoError(t, err)

	err = service.UnregisterUser(user.Token)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not exists")
}

