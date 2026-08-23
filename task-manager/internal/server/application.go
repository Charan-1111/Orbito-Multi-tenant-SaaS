package server

import (
	"task-manager/internal/config"
)

type Application struct {
	config *config.Configuration
}

func NewApplication() (*Application, error) {
	config := &config.Configuration{}
	err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	return &Application{
		config: config,
	}, nil
}

func (app *Application) StartApplication() error {
	appServer := app.SetUpRoutes()

	err := appServer.Listen(app.config.Server.Port)

	return err
}
