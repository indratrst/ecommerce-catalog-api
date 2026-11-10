package middleware

import (
	"strings"

	"ecommerce-catalog-api/pkg/response"
	"ecommerce-catalog-api/utils"

	"github.com/gofiber/fiber/v2"
)

// Protected memverifikasi apakah request membawa Bearer Token JWT yang valid
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return authenticate(c, false)
	}
}

// OptionalProtected mengizinkan request tanpa token dan mengisi user pada
// context ketika token valid diberikan.
func OptionalProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return authenticate(c, true)
	}
}

func authenticate(c *fiber.Ctx, optional bool) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		if optional {
			return c.Next()
		}
		return response.Error(c, fiber.StatusUnauthorized, "Akses ditolak: Token autentikasi tidak ditemukan", nil)
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return response.Error(c, fiber.StatusUnauthorized, "Format header Authorization harus 'Bearer <token>'", nil)
	}

	claims, err := utils.ValidateToken(parts[1])
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, err.Error(), nil)
	}

	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return response.Error(c, fiber.StatusUnauthorized, "Payload token tidak valid: user_id hilang", nil)
	}

	role, _ := claims["role"].(string)
	c.Locals("user_id", userID)
	c.Locals("role", role)

	return c.Next()
}

// RequireRole membatasi akses endpoint hanya untuk role tertentu (misal: "admin")
func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role").(string)
		if !ok || userRole != requiredRole {
			return response.Error(c, fiber.StatusForbidden, "Akses ditolak: Anda tidak memiliki izin untuk sumber daya ini", nil)
		}
		return c.Next()
	}
}
