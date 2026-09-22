package middleware

import (
	"ecommerce-catalog-api/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Header Authorization tidak ditemukan"})
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		userID, err := utils.ValidateToken(tokenString)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": err.Error()})
		}

		// Simpan user_id di dalam Locals Fiber agar bisa diakses di handler selanjutnya
		c.Locals("user_id", userID)
		return c.Next()
	}
}
