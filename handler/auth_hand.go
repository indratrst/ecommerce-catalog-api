package handler

import (
	"my-go-api/domain"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService domain.AuthService
}

func NewAuthHandler(authService domain.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req domain.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	res, err := h.authService.Register(req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Register berhasil", "data": res})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	res, err := h.authService.Login(req)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Login berhasil", "data": res})
}