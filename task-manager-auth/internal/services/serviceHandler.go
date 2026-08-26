package services

import (
	"task-manager-auth/internal/config"
	"task-manager-auth/internal/logging"
	repository "task-manager-auth/internal/store/storeRepo"
)

type Service struct {
	config   *config.Configuration
	log      *logging.Log
	database repository.Repository
}

func NewService(config *config.Configuration, log *logging.Log, repo repository.Repository) *Service {
	return &Service{
		config:   config,
		log:      log,
		database: repo,
	}
}
