package server

import (
	"context"
	"task-manager/internal/config"
	"task-manager/internal/logging"
	"task-manager/internal/store/database"
	repository "task-manager/internal/store/storeRepo"
)

type Application struct {
	config *config.Configuration
	log    *logging.Log
	db     repository.Repository
}

func NewApplication() (*Application, error) {

	log := &logging.Log{}
	log.Initialize()

	config := &config.Configuration{}
	err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	databaseStore := &database.DataBaseStore{}
	err = databaseStore.InitializeDatabaseStore(context.Background(), config.Database)
	if err != nil {
		return nil, err
	}

	return &Application{
		config: config,
		log:    log,
		db:     databaseStore,
	}, nil
}

func (app *Application) StartApplication() error {
	appServer := app.SetUpRoutes()

	err := appServer.Listen(app.config.Server.Port)

	return err
}
