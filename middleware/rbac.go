package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// RequirePermission เหลือไว้เพื่อ compatibility — RBAC ย้ายไปใช้ handlers.RequireAdmin / handlers.RequireManagerOrAbove
func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}
