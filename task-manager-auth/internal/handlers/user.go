package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

func (cfg *ConfigHandler) RegisterUser(c fiber.Ctx) error {
	auth := strings.TrimSpace(c.Get("Authorization"))
	requestId := c.Locals("requestId")
	requestID, _ := requestId.(string)
	if requestID == "" {
		requestID = "unknown"
	}

	logger := cfg.log.WithRequestID(requestID)
	logger.Info().Str("event", "register_user").Msg("register user request started")

	err := cfg.Service.RegisterUser(c.Context(), requestID, auth)
	if err != nil {
		logger.Error().Err(err).Str("event", "register_user").Msg("register user request failed")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to register user",
			"message": err.Error(),
		})
	}

	logger.Info().Str("event", "register_user").Msg("register user request completed")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"code": 0, "message": "User registered succeessfulluy"})
}
