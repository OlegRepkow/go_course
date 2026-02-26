package middlewares

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			log.Printf("[WARN] Auth failed: missing Authorization header, path=%s\n", c.Path())
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			log.Printf("[WARN] Auth failed: invalid token, path=%s, error=%v\n", c.Path(), err)
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}

		claims := token.Claims.(jwt.MapClaims)
		username := claims["username"]
		c.Locals("username", username)
		log.Printf("[INFO] Auth success: user=%s, path=%s\n", username, c.Path())
		return c.Next()
	}
}
