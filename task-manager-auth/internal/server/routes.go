package server

import (
	"task-manager-auth/internal/handlers"
	"task-manager-auth/internal/middlewares"
	"task-manager-auth/internal/services"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func (app *Application) SetUpRoutes() *fiber.App {
	appServer := fiber.New()

	// pre-handler information
	configHandler := handlers.NewConfigHandler(app.config, app.log, app.db)
	configHandler.Service = services.NewService(app.config, app.log, app.db)

	// appServer.Use(recover.New())
	appServer.Use(cors.New())

	apiRoutes := appServer.Group("/auth/v1")
	apiRoutes.Use(middlewares.RequestId)

	apiRoutes.Post("/register", configHandler.RegisterUser)

	apiRoutes.Use(middlewares.Authenticate)

	appServer.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("OK")
	})

	appServer.Get("/live", func(c fiber.Ctx) error {
		return c.SendString("OK")
	})
	return appServer
}

// curl --request POST \
//   --url http://localhost:8000/auth/v1/register \
//   --header 'Accept: */*' \
//   --header 'Accept-Encoding: gzip, deflate, br' \
//   --header 'Authorization: Basic Y2hhcmFuOnBhc3N3b3Jk' \
//   --header 'Connection: keep-alive' \
//   --header 'User-Agent: EchoapiRuntime/1.1.0'
