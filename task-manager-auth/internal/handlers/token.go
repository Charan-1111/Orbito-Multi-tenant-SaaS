package handlers

import (
	"task-manager-auth/internal/constants"

	"github.com/gofiber/fiber/v3"
)

type ValidateReponse struct {
	UserId string `json:"userId"`
	Email  string `json:"email"`
}

func (cfg *ConfigHandler) ValidateToken(c fiber.Ctx) error {
	authToken := c.Get("Authorization")
	if authToken == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": 1, "message": "Token not provided"})
	}

	claims, err := cfg.Service.ValidateToken(authToken)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": 1, "message": constants.ErrInvalidToken})
	}

	validateResponse := ValidateReponse{
		UserId: claims.UserId,
		Email:  claims.Email,
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"code": 0, "message": "Token validation successful", "userInfo": validateResponse})
}

func (cfg *ConfigHandler) RotateToken(c fiber.Ctx) error {
	refreshToken := c.Get("Authorization")
	if refreshToken == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": 1, "message": constants.ErrInvalidRefreshToken})
	}

	accessToken, refreshToken, err := cfg.Service.RotateTokens(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": 1, "message": constants.ErrInvalidRefreshToken})
	}

	c.Set("X-Access-Token", accessToken)
	c.Set("X-Refresh-Token", refreshToken)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"code": 0, "message": "Tokens rotated successfully"})
}
