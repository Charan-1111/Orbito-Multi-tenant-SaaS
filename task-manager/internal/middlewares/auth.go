package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

func Authenticate(c fiber.Ctx) error {
	authToken := c.Get("Authorization")

	authToken = strings.Split(authToken, " ")[1]

	c.Set("token", authToken)

	return c.Next()
}
