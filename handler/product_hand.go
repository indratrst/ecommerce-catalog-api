package handler

import (
	"ecommerce-catalog-api/domain"
	"ecommerce-catalog-api/pkg/response"

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
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	res, err := h.service.CreateProduct(c.UserContext(), req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusCreated, "Produk berhasil dibuat", res)
}
func (h *ProductHandler) GetProducts(c *fiber.Ctx) error {
	var param domain.ProductQueryParam
	if err := c.QueryParser(&param); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid query parameters", err.Error())
	}

	res, err := h.service.GetProducts(c.UserContext(), param)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil daftar produk", res)
}

func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	res, err := h.service.GetProductByID(c.UserContext(), id)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "Produk tidak ditemukan", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil detail produk", res)
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var req domain.ProductRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	res, err := h.service.UpdateProduct(c.UserContext(), id, req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Produk berhasil diperbarui", res)
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.DeleteProduct(c.UserContext(), id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Produk berhasil dihapus", nil)
}
