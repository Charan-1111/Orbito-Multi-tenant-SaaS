package handlers

import (
	"task-manager-auth/internal/constants"

	"github.com/gofiber/fiber/v3"
)

func (cfg *ConfigHandler) ValidateToken(c fiber.Ctx) error {
	authToken := c.Get("Authorization")
	if authToken == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": 1, "message": "Token not provided"})
	}

	claims, err := cfg.Service.ValidateToken(authToken)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": 1, "message": constants.ErrInvalidToken})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"code": 0, "message": "Token validation successful", "claims": claims})
}
