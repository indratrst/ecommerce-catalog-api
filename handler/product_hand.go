package handler

import (
	"fmt"
	"my-go-api/domain"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	service domain.ProductService
}

func NewProductHandler(service domain.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req domain.ProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	// Logging sederhana untuk debug di terminal docker
	fmt.Printf("Request Payload: %+v\n", req)

	product, err := h.service.CreateProduct(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"data": product})
}

func (h *ProductHandler) GetProducts(c *fiber.Ctx) error {
	var param domain.ProductQueryParam

	// Parse query params seperti ?page=1&limit=5&search=keyboard
	if err := c.QueryParser(&param); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Query parameter tidak valid"})
	}

	result, err := h.service.GetAllProducts(param)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	product, err := h.service.GetProductByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Produk tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"data": product})
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var req domain.ProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	product, err := h.service.UpdateProduct(id, req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": product})
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.DeleteProduct(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Produk berhasil dihapus"})
}
