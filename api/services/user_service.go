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

	count, err := service.Collection.CountDocuments(ctx, bson.M{"_id": token})
	if err != nil {
		return err
	}
	
	if count <= 0 {
		return fmt.Errorf("user with token %s does not exists", token)
	}

	_, err = service.Collection.DeleteOne(ctx, bson.M{"_id": token})
	return err
}
