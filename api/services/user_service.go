package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Data struct {
	ApiToken string `bson:"_id" json:"id"`
}

type UserService struct {
	Collection *mongo.Collection
}

const TIMEOUT = 5 * time.Second

func (service *UserService) RegisterUser() (Data ,error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	rawID := fmt.Sprintf("%s-%s", uuid.New(), time.Now().String())
	identifier := sha256.Sum256([]byte(rawID))

	apiToken := hex.EncodeToString(identifier[:])

	var data Data
	count, err := service.Collection.CountDocuments(ctx, bson.M{"_id": apiToken})
	if err != nil {
		return data, err
	}

	if count > 0 {
		return data, fmt.Errorf("This token already exists in the system")
	}

	data = Data {
		ApiToken: apiToken,
	}

	_, err = service.Collection.InsertOne(ctx, data)
	return data, err
}

func (service *UserService) UnregisterUser(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	tokenHash := sha256.Sum256([]byte(token))
	tokenStr := hex.EncodeToString(tokenHash[:])
	count, err := service.Collection.CountDocuments(ctx, bson.M{"_id": tokenStr})
	if err != nil {
		return err
	}

	if count <= 0 {
		return fmt.Errorf("document with token %s does not exists", token)
	}

	_, err = service.Collection.DeleteOne(ctx, bson.M{"_id": tokenStr})
	return err
}
