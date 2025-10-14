package services

import (
	"context"
	"time"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RawUser struct {
	Owner string `json:"owner"`
	Subject string `json:"subject"`
	Subscribers []string `json:"subscribers,omitempty"`
}

type User struct {
	Owner string `bson:"_id" json:"id"`
	Token string `bson:"token" json:"token"`
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

	count, err := service.Collection.CountDocuments(ctx, bson.M{"_id": user.Owner})
	if err != nil {
		return err
	}
	
	if count > 0 {
		return fmt.Errorf("user with owner %s already exists", user.Owner)
	}

	_, err = service.Collection.InsertOne(ctx, user)
	return err
}

func (service *UserService) UnregisterUser(owner string) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	count, err := service.Collection.CountDocuments(ctx, bson.M{"_id": owner})
	if err != nil {
		return err
	}
	
	if count <= 0 {
		return fmt.Errorf("user with owner %s does not exists", owner)
	}

	_, err = service.Collection.DeleteOne(ctx, bson.M{"_id": owner})
	return err
}
