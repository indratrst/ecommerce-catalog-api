package handler

import (
	"ecommerce-catalog-api/domain"
	"ecommerce-catalog-api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	categoryService domain.CategoryService
}

func NewCategoryHandler(categoryService domain.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}
func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req domain.CategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	res, err := h.categoryService.CreateCategory(c.UserContext(), req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusCreated, "Kategori berhasil dibuat", res)
}

func (h *CategoryHandler) GetCategories(c *fiber.Ctx) error {
	res, err := h.categoryService.GetCategories(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil data kategori", res)
}

func (h *CategoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	res, err := h.categoryService.GetCategoryByID(c.UserContext(), c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "Kategori tidak ditemukan", err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Berhasil mengambil detail kategori", res)
}

func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	var req domain.CategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	res, err := h.categoryService.UpdateCategory(c.UserContext(), c.Params("id"), req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, fiber.StatusOK, "Kategori berhasil diperbarui", res)
}

func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	if err := h.categoryService.DeleteCategory(c.UserContext(), c.Params("id")); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, fiber.StatusOK, "Kategori berhasil dihapus", nil)
}
