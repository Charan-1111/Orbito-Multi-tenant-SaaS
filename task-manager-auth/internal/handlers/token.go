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
	requestId, _ := c.Locals("requestId").(string)
	if requestId == "" {
		requestId = "unknown"
	}

	logger := cfg.log.WithRequestID(requestId)
	logger.Info().Str("event", "validate_token").Msg("token validation request started")

	if authToken == "" {
		logger.Warn().Str("event", "validate_token").Msg("token validation request missing authorization header")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": 1, "message": "Token not provided"})
	}

	claims, err := cfg.Service.ValidateToken(authToken)
	if err != nil {
		logger.Error().Err(err).Str("event", "validate_token").Msg("token validation failed")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": 1, "message": constants.ErrInvalidToken})
	}

	validateResponse := ValidateReponse{
		UserId: claims.UserId,
		Email:  claims.Email,
	}
	logger.Info().Str("event", "validate_token").Msg("token validation request completed")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"code": 0, "message": "Token validation successful", "userInfo": validateResponse})
}

func (cfg *ConfigHandler) RotateToken(c fiber.Ctx) error {
	refreshToken := c.Get("Authorization")
	requestId, _ := c.Locals("requestId").(string)
	if requestId == "" {
		requestId = "unknown"
	}

	logger := cfg.log.WithRequestID(requestId)
	logger.Info().Str("event", "rotate_token").Msg("token rotation request started")

	if refreshToken == "" {
		logger.Warn().Str("event", "rotate_token").Msg("token rotation request missing authorization header")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": 1, "message": constants.ErrInvalidRefreshToken})
	}

	accessToken, refreshToken, err := cfg.Service.RotateTokens(refreshToken)
	if err != nil {
		logger.Error().Err(err).Str("event", "rotate_token").Msg("token rotation failed")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": 1, "message": constants.ErrInvalidRefreshToken})
	}

	c.Set("X-Access-Token", accessToken)
	c.Set("X-Refresh-Token", refreshToken)

	logger.Info().Str("event", "rotate_token").Msg("token rotation request completed")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"code": 0, "message": "Tokens rotated successfully"})
}
