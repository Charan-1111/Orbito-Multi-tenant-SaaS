package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

func Authenticate(c fiber.Ctx) error {
	authHeader := strings.TrimSpace(c.Get("Authorization"))
	parts := strings.Fields(authHeader)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return fiber.ErrUnauthorized
	}

	c.Locals("token", parts[1])

	// TODO : we need to implemet the validate logic here

	return c.Next()
}
