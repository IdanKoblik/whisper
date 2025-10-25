package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
	"whisper-api/config"
	"whisper-api/db"

	"github.com/google/uuid"
)

type User struct {
	ApiToken string `bson:"_id" json:"id"`
}

type RegisterResponse struct {
	ApiToken string `json:"api_token"`
}

func RegisterUser(cfg *config.Config) (RegisterResponse, error) {
	rawID := fmt.Sprintf("%s-%s", uuid.New(), time.Now().String())
	identifier := sha256.Sum256([]byte(rawID))

	apiToken := hex.EncodeToString(identifier[:])
	data := User{
		ApiToken: apiToken,
	}

	var response RegisterResponse
	err := db.InsertData(cfg, data, apiToken)
	if err != nil {
		return response, err
	}

	response = RegisterResponse{ApiToken: apiToken}
	return response, nil
}

func RemoveUser(cfg *config.Config, apiToken string) error {
	rawID := fmt.Sprintf("%s-%s", uuid.New(), time.Now().String())
	identifier := sha256.Sum256([]byte(rawID))

	hashedApiToken := hex.EncodeToString(identifier[:])
	err := db.DeleteData(cfg, hashedApiToken)
	if err != nil {
		return err
	}

	err = db.DeleteData(cfg, apiToken)
	return err
}
