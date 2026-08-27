package handlers

import (
	"task-manager-auth/internal/config"
	"task-manager-auth/internal/logging"
	"task-manager-auth/internal/services"
	"task-manager-auth/internal/store/database"
)

type ConfigHandler struct {
	config   *config.Configuration
	log      *logging.Log
	database database.Repository
	Service  *services.Service
}

func NewConfigHandler(config *config.Configuration, log *logging.Log, database database.Repository) *ConfigHandler {
	return &ConfigHandler{
		config:   config,
		log:      log,
		database: database,
	}
}
