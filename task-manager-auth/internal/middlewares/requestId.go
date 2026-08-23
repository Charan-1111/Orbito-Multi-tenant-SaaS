package middlewares

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func RequestId(c fiber.Ctx) error {
	requestID := c.Get("X-Request-Id")
	if requestID == "" {
		requestID = uuid.NewString()
	}
	c.Set("X-Request-Id", requestID)
	c.Locals("requestid", requestID)

	return c.Next()
}
