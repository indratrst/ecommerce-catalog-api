package handler

import (
	"ecommerce-catalog-api/domain"
	"ecommerce-catalog-api/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type AuthHandler struct {
	authService domain.AuthService
}

func NewAuthHandler(authService domain.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req domain.RegisterRequest

	// 1. Parse Request Body
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Format payload JSON tidak valid", nil)
	}

	// 2. Validasi Struct Field
	if err := validate.Struct(&req); err != nil {
		validationErrors := formatValidationErrors(err)
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validasi gagal", validationErrors)
	}

	// 3. Panggil Service/Usecase
	res, err := h.authService.Register(c.Context(), req)
	if err != nil {
		// Jika email sudah terdaftar
		if err.Error() == "email already registered" {
			return response.Error(c, fiber.StatusConflict, "Email sudah terdaftar", nil)
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	// 4. Return Success Response
	return response.Success(c, fiber.StatusCreated, "Registrasi berhasil", res)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest

	// 1. Parse Request Body
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Format payload JSON tidak valid", nil)
	}

	// 2. Validasi Struct Field
	if err := validate.Struct(&req); err != nil {
		validationErrors := formatValidationErrors(err)
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validasi gagal", validationErrors)
	}

	// 3. Panggil Service/Usecase
	res, err := h.authService.Login(c.Context(), req)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "Email atau password salah", nil)
	}

	// 4. Return Success Response
	return response.Success(c, fiber.StatusOK, "Login berhasil", res)
}

// Helper untuk memformat error validator menjadi array JSON yang rapi
func formatValidationErrors(err error) []map[string]string {
	var errors []map[string]string
	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, map[string]string{
			"field":   err.Field(),
			"message": err.Tag() + " validation failed",
		})
	}
	return errors
}
