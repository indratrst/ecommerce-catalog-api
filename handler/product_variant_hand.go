package handler

import (
	"ecommerce-catalog-api/domain"
	"ecommerce-catalog-api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type ProductVariantHandler struct {
	service domain.ProductVariantService
}

func NewProductVariantHandler(service domain.ProductVariantService) *ProductVariantHandler {
	return &ProductVariantHandler{service: service}
}

func (h *ProductVariantHandler) CreateProductVariant(c *fiber.Ctx) error {
	productID := c.Params("productId")
	var req domain.ProductVariantRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	res, err := h.service.CreateProductVariant(c.UserContext(), productID, req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusCreated, "Variant produk berhasil dibuat", res)
}

func (h *ProductVariantHandler) GetProductVariantsByProductID(c *fiber.Ctx) error {
	productID := c.Params("productId")
	res, err := h.service.GetProductVariantsByProductID(c.UserContext(), productID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil daftar variant produk", res)
}

func (h *ProductVariantHandler) GetProductVariantByID(c *fiber.Ctx) error {
	id := c.Params("id")
	res, err := h.service.GetProductVariantByID(c.UserContext(), id)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "Variant produk tidak ditemukan", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil detail variant produk", res)
}

func (h *ProductVariantHandler) UpdateProductVariant(c *fiber.Ctx) error {
	id := c.Params("id")
	var req domain.ProductVariantRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	res, err := h.service.UpdateProductVariant(c.UserContext(), id, req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Variant produk berhasil diperbarui", res)
}

func (h *ProductVariantHandler) DeleteProductVariant(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.DeleteProductVariant(c.UserContext(), id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Variant produk berhasil dihapus", nil)
}
