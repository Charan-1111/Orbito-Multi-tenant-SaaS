package middlewares

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestAuthenticateRejectsMalformedAuthorization(t *testing.T) {
	app := fiber.New()
	app.Get("/", Authenticate, func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	for _, header := range []string{"", "Token abc", "Bearer", "Bearer   ", "Bearer abc extra"} {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", header)

		res, err := app.Test(req)
		if err != nil {
			t.Fatalf("app.Test() error = %v", err)
		}
		if res.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("header %q status = %d, want %d", header, res.StatusCode, fiber.StatusUnauthorized)
		}
	}
}

func TestAuthenticateStoresBearerToken(t *testing.T) {
	app := fiber.New()
	app.Get("/", Authenticate, func(c fiber.Ctx) error {
		if token := c.Locals("token"); token != "abc123" {
			t.Errorf("token = %q, want %q", token, "abc123")
		}
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "bearer abc123")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("status = %d, want %d", res.StatusCode, fiber.StatusOK)
	}
}
