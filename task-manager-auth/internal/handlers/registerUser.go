package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

func (cfg *ConfigHandler) RegisterUser(c fiber.Ctx) error {

	auth := strings.TrimSpace(c.Get("Authorization"))

	requestId := c.Locals("requestId")

	err := cfg.Service.RegisterUser(c.Context(), requestId.(string), auth)
	if err != nil {
		// need to have proper error messagew
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to register user",
			"message": err.Error(),
		})
	}

	return nil
}
