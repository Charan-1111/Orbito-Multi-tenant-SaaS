package server

import (
	"context"
	"os"
	"sync"
	"task-manager/internal/config"
	"task-manager/internal/logging"
	"task-manager/internal/store/database"
	repository "task-manager/internal/store/storeRepo"
	"time"

	"github.com/gofiber/fiber/v3"
)

type Application struct {
	config *config.Configuration
	log    *logging.Log
	db     repository.Repository
	server *fiber.App
	once   sync.Once
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

func (app *Application) StartApplication(shutdown <-chan os.Signal) error {
	app.server = app.SetUpRoutes()
	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- app.server.Listen(app.config.Server.Port)
	}()

	select {
	case err := <-serverErrors:
		app.closeResources()
		return err
	case <-shutdown:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		shutdownErr := app.server.ShutdownWithContext(ctx)
		listenErr := <-serverErrors
		app.closeResources()
		if shutdownErr != nil {
			return shutdownErr
		}
		return listenErr
	}
}

func (app *Application) closeResources() {
	app.once.Do(func() {
		if app.db != nil {
			app.db.Close()
		}
		if app.log != nil {
			app.log.Close()
		}
	})
}
