package middlewares

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func TestAuthMiddleware(t *testing.T) {
	secret := "testsecret"
	app := fiber.New()

	app.Use(Auth(secret))
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	t.Run("No token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		resp, _ := app.Test(req)
		if resp.StatusCode != 401 {
			t.Fatalf("Expected 401, got %d", resp.StatusCode)
		}
	})

	t.Run("Valid token", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": "testuser",
			"exp":      time.Now().Add(time.Hour).Unix(),
		})
		tokenStr, _ := token.SignedString([]byte(secret))

		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+tokenStr)
		resp, _ := app.Test(req)
		if resp.StatusCode != 200 {
			t.Fatalf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Invalid token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid")
		resp, _ := app.Test(req)
		if resp.StatusCode != 401 {
			t.Fatalf("Expected 401, got %d", resp.StatusCode)
		}
	})
}
