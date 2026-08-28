package services

import (
	"task-manager-auth/internal/config"
	"task-manager-auth/internal/logging"
	"task-manager-auth/internal/store/database"
	"task-manager-auth/internal/token"
)

type Service struct {
	config       *config.Configuration
	log          *logging.Log
	database     database.Repository
	tokenService *token.TokenService
}

func NewService(config *config.Configuration, log *logging.Log, repo database.Repository, tokenService *token.TokenService) *Service {
	return &Service{
		config:       config,
		log:          log,
		database:     repo,
		tokenService: tokenService,
	}
}
