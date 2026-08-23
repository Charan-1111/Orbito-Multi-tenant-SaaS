package middlewares

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestRequestIdPreservesIncomingID(t *testing.T) {
	app := fiber.New()
	app.Get("/", RequestId, func(c fiber.Ctx) error {
		return c.SendString(c.Get("X-Request-Id"))
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Request-Id", "incoming-id")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.Header.Get("X-Request-Id") != "incoming-id" {
		t.Fatalf("response request ID = %q, want %q", res.Header.Get("X-Request-Id"), "incoming-id")
	}
}

func TestRequestIdGeneratesMissingID(t *testing.T) {
	app := fiber.New()
	app.Get("/", RequestId, func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	res, err := app.Test(httptest.NewRequest("GET", "/", nil))
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.Header.Get("X-Request-Id") == "" {
		t.Fatal("response request ID is empty")
	}
}
