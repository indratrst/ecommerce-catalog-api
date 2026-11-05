package handler

import (
	"ecommerce-catalog-api/domain"
	"ecommerce-catalog-api/pkg/response"
	"log"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	service domain.OrderService
}

func NewOrderHandler(service domain.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var req domain.OrderRequest
	log.Println("===> 1. Hit Endpoint Checkout Order")
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ERROR] Body Parser: %v", err)
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if userID, ok := c.Locals("user_id").(string); ok && userID != "" {
		req.UserID = userID
	}
	log.Printf("===> 3. Payload Request: %+v", req)

	res, err := h.service.CreateOrder(c.UserContext(), req)
	if err != nil {
		log.Printf("[ERROR] Service CreateOrder: %v", err)
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	log.Println("===> 4. Berhasil Simpan ke Database!")
	return response.Success(c, fiber.StatusCreated, "Order created successfully", res)
}

func (h *OrderHandler) GetOrders(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Error(c, fiber.StatusUnauthorized, "User not authenticated", nil)
	}

	res, err := h.service.GetOrders(c.UserContext(), userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Orders retrieved successfully", res)
}

func (h *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	orderID := c.Params("id")
	res, err := h.service.GetOrderByID(c.UserContext(), orderID)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "Order not found", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Order retrieved successfully", res)
}

func (h *OrderHandler) UpdateShipment(c *fiber.Ctx) error {
	orderID := c.Params("id")
	var req domain.OrderShipmentRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	res, err := h.service.UpdateShipment(c.UserContext(), orderID, req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Shipment updated successfully", res)
}

func (h *OrderHandler) ProcessPaymentWebhook(c *fiber.Ctx) error {
	payload := make(map[string]interface{})
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid webhook payload", err.Error())
	}

	if err := h.service.ProcessPaymentWebhook(c.UserContext(), payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Webhook processed successfully", nil)
}
