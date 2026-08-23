package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
)

func RequestId(c fiber.Ctx) error {
	requestId := c.Get("X-Request-Id")

	fmt.Println("Request ID:", requestId)

	if requestId == "" {
		c.Set("X-Request-Id", c.Locals("requestid").(string))
	}

	return c.Next()
}
