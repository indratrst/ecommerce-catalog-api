package handler

import (
	"ecommerce-catalog-api/domain"
	"ecommerce-catalog-api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type AdminUserHandler struct {
	service domain.AdminUserService
}

func NewAdminUserHandler(service domain.AdminUserService) *AdminUserHandler {
	return &AdminUserHandler{service: service}
}

func (h *AdminUserHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.service.GetUsers(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, fiber.StatusOK, "Berhasil mengambil data pengguna", users)
}

func (h *AdminUserHandler) CreateUser(c *fiber.Ctx) error {
	var req domain.AdminUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	user, err := h.service.CreateUser(c.UserContext(), req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	return response.Success(c, fiber.StatusCreated, "Pengguna berhasil dibuat", user)
}

func (h *AdminUserHandler) UpdateUser(c *fiber.Ctx) error {
	var req domain.AdminUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	user, err := h.service.UpdateUser(c.UserContext(), c.Params("id"), req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	return response.Success(c, fiber.StatusOK, "Pengguna berhasil diperbarui", user)
}

func (h *AdminUserHandler) DeleteUser(c *fiber.Ctx) error {
	currentUserID, _ := c.Locals("user_id").(string)
	if c.Params("id") == currentUserID {
		return response.Error(c, fiber.StatusBadRequest, "Akun yang sedang digunakan tidak dapat dihapus", nil)
	}
	if err := h.service.DeleteUser(c.UserContext(), c.Params("id")); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, fiber.StatusOK, "Pengguna berhasil dihapus", nil)
}
