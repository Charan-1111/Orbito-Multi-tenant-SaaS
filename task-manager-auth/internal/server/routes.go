package server

import (
	"task-manager-auth/internal/handlers"
	"task-manager-auth/internal/middlewares"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func (app *Application) SetUpRoutes() *fiber.App {
	appServer := fiber.New()

	// appServer.Use(recover.New())
	appServer.Use(cors.New())

	apiRoutes := appServer.Group("/auth/v1")
	apiRoutes.Use(middlewares.RequestId)
	apiRoutes.Use(middlewares.Authenticate)

	// pre-handler information
	configHandler := handlers.NewConfigHandler(app.config, app.log, app.db)

	// Auth routes
	apiRoutes.Post("/register", configHandler.RegisterUser)

	appServer.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("OK")
	})

	appServer.Get("/live", func(c fiber.Ctx) error {
		return c.SendString("OK")
	})
	return appServer
}
