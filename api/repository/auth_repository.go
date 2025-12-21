package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"whisper-api/models"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepository struct {
	Col      *mongo.Collection
	Rdb      *redis.Client
	CacheTTL time.Duration
}

var (
	REDIS_AUTH_PREFIX = "auth:"
)

func NewAuthRepository(mongoCli *mongo.Client, redisCli *redis.Client, dbName, colName string) *AuthRepository {
	return &AuthRepository{
		Rdb:      redisCli,
		Col:      mongoCli.Database(dbName).Collection(colName),
		CacheTTL: 24 * time.Hour,
	}
}

func (r *AuthRepository) getCacheKey(token string) string {
	return REDIS_AUTH_PREFIX + token
}

func (r *AuthRepository) CreateToken(ctx context.Context, data *models.AuthModel) error {
	_, err := r.Col.InsertOne(ctx, data)
	return err
}

func (r *AuthRepository) ValidateToken(ctx context.Context, token string) (*models.AuthModel, error) {
	key := r.getCacheKey(token)
	val, err := r.Rdb.Get(ctx, key).Result()
	if err == nil {
		var data models.AuthModel
		json.Unmarshal([]byte(val), &data)
		return &data, nil
	}

	var data models.AuthModel
	err = r.Col.FindOne(ctx, bson.M{"_id": token}).Decode(&data)
	if err != nil {
		return nil, err
	}

	b, _ := json.Marshal(data)
	r.Rdb.Set(ctx, key, b, r.CacheTTL)
	return &data, nil
}

func (r *AuthRepository) ValidateDeviceID(ctx context.Context, tokenHash, deviceID string) (bool, error) {
	data, err := r.ValidateToken(ctx, tokenHash)
	if err != nil {
		return false, err
	}

	for _, d := range data.Devices {
		if d == deviceID {
			return true, nil
		}
	}
	return false, nil
}

func (r *AuthRepository) RemoveDeviceID(ctx context.Context, tokenHash, deviceID string) error {
	filter := bson.M{"_id": tokenHash}
	update := bson.M{"$pull": bson.M{"devices": deviceID}}

	_, err := r.Col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return r.Rdb.Del(ctx, r.getCacheKey(tokenHash)).Err()
}

func (r *AuthRepository) RemoveToken(ctx context.Context, tokenHash string) error {
	res, err := r.Col.DeleteOne(ctx, bson.M{"_id": tokenHash})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return fmt.Errorf("token not found")
	}

	return r.Rdb.Del(ctx, r.getCacheKey(tokenHash)).Err()
}

func (r *AuthRepository) AddDeviceID(ctx context.Context, tokenHash, deviceID string) error {
	data, err := r.ValidateToken(ctx, tokenHash)
	if err != nil {
		return err
	}

	for _, d := range data.Devices {
		if d == deviceID {
			return nil
		}
	}

	filter := bson.M{"_id": tokenHash}
	update := bson.M{"$addToSet": bson.M{"devices": deviceID}}

	_, err = r.Col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	data.Devices = append(data.Devices, deviceID)
	b, _ := json.Marshal(data)
	err = r.Rdb.Set(ctx, r.getCacheKey(tokenHash), b, r.CacheTTL).Err()
	return err
}
