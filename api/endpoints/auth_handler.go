package endpoints

import "whisper-api/repository"

type AuthHandler struct {
	Repo *repository.AuthRepository
}

func NewAuthHandler(repo *repository.AuthRepository) *AuthHandler {
	return &AuthHandler{Repo: repo}
}
