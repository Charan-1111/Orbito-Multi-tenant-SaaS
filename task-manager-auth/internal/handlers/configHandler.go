package handlers

import (
	"task-manager-auth/internal/config"
	"task-manager-auth/internal/logging"
	repository "task-manager-auth/internal/store/storeRepo"
)

type ConfigHandler struct {
	config   *config.Configuration
	log      *logging.Log
	database repository.Repository
}


func NewConfigHandler(config *config.Configuration, log *logging.Log, database repository.Repository) *ConfigHandler {
	return &ConfigHandler{
		config:   config,
		log: log,
		database: database,
	}
}