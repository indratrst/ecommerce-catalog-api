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
		authHeader := c.Get("Authorization")

		// 1. Cek keberadaan Header Authorization
		if authHeader == "" {
			return response.Error(c, fiber.StatusUnauthorized, "Akses ditolak: Token autentikasi tidak ditemukan", nil)
		}

		// 2. Memastikan format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.Error(c, fiber.StatusUnauthorized, "Format header Authorization harus 'Bearer <token>'", nil)
		}

		tokenString := parts[1]

		// 3. Validasi Token JWT
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, err.Error(), nil)
		}

		// 4. Ekstrak user_id dan role dari claims
		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			return response.Error(c, fiber.StatusUnauthorized, "Payload token tidak valid: user_id hilang", nil)
		}

		role, _ := claims["role"].(string)

		// 5. Simpan ke c.Locals agar bisa diakses di Handler/Controller (misal: c.Locals("user_id"))
		c.Locals("user_id", userID)
		c.Locals("role", role)

		return c.Next()
	}
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