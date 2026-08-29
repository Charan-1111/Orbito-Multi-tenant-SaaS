package server

import (
	"context"
	"os"
	"sync"
	"task-manager-auth/internal/config"
	"task-manager-auth/internal/constants"
	"task-manager-auth/internal/logging"
	"task-manager-auth/internal/store/database"
	"task-manager-auth/internal/token"
	"time"

	"github.com/gofiber/fiber/v3"
)

type Application struct {
	config       *config.Configuration
	log          *logging.Log
	db           database.Repository
	server       *fiber.App
	tokenService *token.TokenService
	once         sync.Once
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
	databaseStore.Log = log
	err = databaseStore.InitializeDatabaseStore(context.Background(), config.Queries)
	if err != nil {
		return nil, err
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, constants.ErrMissingJwtSecret
	}

	tokenService, err := token.NewTokenService(secret, constants.TokenIssuer)
	if err != nil {
		return nil, err
	}

	return &Application{
		config:       config,
		log:          log,
		db:           databaseStore,
		tokenService: tokenService,
	}, nil
}

func (app *Application) StartApplication(shutdown <-chan os.Signal) error {
	app.PerformNecessaryOperations()
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

func (app *Application) PerformNecessaryOperations() {
	err := app.db.CreateTables(context.Background())
	if err != nil {
		app.log.Log.Error().Err(err).Msg("Error initializing the tables")
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
