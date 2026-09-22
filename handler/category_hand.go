package handler

import (
	"my-go-api/domain"

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
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	category, err := h.categoryService.CreateCategory(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Kategori berhasil dibuat", "data": category})
}

func (h *CategoryHandler) GetCategories(c *fiber.Ctx) error {
	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": categories})
}
