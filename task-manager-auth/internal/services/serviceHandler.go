package services

import (
	"task-manager-auth/internal/config"
	"task-manager-auth/internal/logging"
	"task-manager-auth/internal/store/database"
)

type Service struct {
	config   *config.Configuration
	log      *logging.Log
	database database.Repository
}

func NewService(config *config.Configuration, log *logging.Log, repo database.Repository) *Service {
	return &Service{
		config:   config,
		log:      log,
		database: repo,
	}
}
