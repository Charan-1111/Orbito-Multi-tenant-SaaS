package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

func (h *ConfigHandler) UserLogin(c fiber.Ctx) error {
	auth := strings.TrimSpace(c.Get("Authorization"))

	requestId := c.Locals("requestId").(string)

	err := h.Service.UserLogin(c.Context(), requestId, auth)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message" : err.Error()})
	}
	
	/*
		TODO:
		After successful login, we need to generate J-Tokens, so that the subsequent api calls
		will use those tokens for the validation
	*/

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message" : "Login Successfull"})
}
