package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
	"whisper-api/config"
	"whisper-api/db"
	"whisper-api/utils"

	"github.com/google/uuid"
)

type User struct {
	ApiToken string `bson:"_id" json:"id"`
}

type RegisterResponse struct {
	ApiToken string `json:"api_token"`
}

func RegisterUser(cfg *config.Config) (RegisterResponse, error) {
	rawToken := fmt.Sprintf("%s-%s", uuid.New(), time.Now().String())
	hash := sha256.Sum256([]byte(rawToken))
	token := hex.EncodeToString(hash[:])

	data := User{
		ApiToken: utils.HashToken(token),
	}

	var response RegisterResponse
	err := db.InsertData(cfg, data, data.ApiToken)
	if err != nil {
		return response, err
	}

	response = RegisterResponse{ApiToken: token}
	return response, nil
}

func RemoveUser(cfg *config.Config, rawToken string) error {
	return db.DeleteData(cfg, utils.HashToken(rawToken))
}
