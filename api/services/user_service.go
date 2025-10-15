package services

import (
	"context"
	"time"
	"fmt"
	"crypto/sha256"
	"encoding/hex"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RawUser struct {
	Owner string `json:"owner"`
	Subject string `json:"subject"`
	Subscribers []string `json:"subscribers,omitempty"`
}

type User struct {
	Token string `bson:"_id" json:"id"`
	Owner string `bson:"owner" json:"owner"`
	Subject string `bson:"subject" json:"subject"`
	Subscribers []string `bson:"subscribers" json:"subscribers"`
}

type UserService struct {
	Collection *mongo.Collection
}

const TIMEOUT = 5 * time.Second

func (service *UserService) RegisterUser(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	tokenHash := sha256.Sum256([]byte(user.Token))
	user.Token = hex.EncodeToString(tokenHash[:])
	count, err := service.Collection.CountDocuments(ctx, bson.M{"_id": user.Token})
	if err != nil {
		return err
	}
	
	if count > 0 {
		return fmt.Errorf("This token already exists in the system")
	}

	_, err = service.Collection.InsertOne(ctx, user)
	return err
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
		return fmt.Errorf("user with token %s does not exists", token)
	}

	_, err = service.Collection.DeleteOne(ctx, bson.M{"_id": tokenStr})
	return err
}
