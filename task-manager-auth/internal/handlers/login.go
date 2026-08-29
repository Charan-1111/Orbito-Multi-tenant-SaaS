package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

func (h *ConfigHandler) UserLogin(c fiber.Ctx) error {
	auth := strings.TrimSpace(c.Get("Authorization"))

	requestId, _ := c.Locals("requestId").(string)
	if requestId == "" {
		requestId = "unknown"
	}

	logger := h.log.WithRequestID(requestId)
	logger.Info().Str("event", "login").Msg("login request started")

	accessToken, refreshToken, err := h.Service.UserLogin(c.Context(), requestId, auth)
	if err != nil {
		logger.Error().Err(err).Str("event", "login").Msg("login request failed")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": err.Error()})
	}

	// injecting the tokens into response headers
	c.Set("X-Access-Token", accessToken)
	c.Set("X-Refresh-Token", refreshToken)

	logger.Info().Str("event", "login").Msg("login request completed")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Login Successfull"})
}
