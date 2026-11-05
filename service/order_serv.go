package service

import (
	"context"
	"ecommerce-catalog-api/domain"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type orderService struct {
	orderRepo domain.OrderRepository
}

func nullableString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func NewOrderService(orderRepo domain.OrderRepository) domain.OrderService {
	return &orderService{orderRepo: orderRepo}
}

func (s *orderService) CreateOrder(ctx context.Context, req domain.OrderRequest) (*domain.OrderResponse, error) {
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("items are required")
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil || userID == uuid.Nil {
		return nil, fmt.Errorf("user is required")
	}

	shippingMethod := domain.ShippingMethodPickupStore
	if req.DeliveryType == "courier" || req.ShippingMethod == "SHIPPING" || req.ShippingMethod == "shipping" {
		shippingMethod = domain.ShippingMethodShipping
	}
	if req.DeliveryType == "pickup" || req.DeliveryType == "PICKUP_STORE" {
		shippingMethod = domain.ShippingMethodPickupStore
	}

	orderID := uuid.New()
	orderNumber := fmt.Sprintf("INV/%s/%04d", time.Now().Format("20060102"), time.Now().UnixNano()%10000)
	shippingAddress := req.ShippingAddress
	if shippingAddress == "" {
		shippingAddress = req.StoreLocation
	}

	var orderItems []*domain.OrderItem
	var totalAmount int = 0

	// 1. Loop items & hitung totalAmount TERLEBIH DAHULU
	for _, itemReq := range req.Items {
		if itemReq.Quantity <= 0 {
			continue
		}

		// Fix: Ambil ProductVariantID, jika kosong baru fallback ke VariantID
		variantIDStr := itemReq.ProductVariantID
		if variantIDStr == "" {
			variantIDStr = itemReq.VariantID
		}

		// Jika keduanya tetap kosong, lempar error jelas
		if variantIDStr == "" {
			return nil, fmt.Errorf("productVariantId / variantId is required")
		}

		variantID, err := uuid.Parse(variantIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid productVariantId: %s", variantIDStr)
		}

		itemPrice := itemReq.PriceAtPurchase
		itemSubtotal := itemPrice * itemReq.Quantity
		totalAmount += itemSubtotal

		item := &domain.OrderItem{
			ID:               uuid.New(),
			OrderID:          orderID,
			ProductVariantID: variantID,
			Quantity:         itemReq.Quantity,
			PriceAtPurchase:  itemPrice,
			VariantSnapshot: map[string]interface{}{
				"productVariantId": variantIDStr,
				"quantity":         itemReq.Quantity,
				"priceAtPurchase":  itemPrice,
			},
		}
		orderItems = append(orderItems, item)
	}

	if len(orderItems) == 0 {
		return nil, fmt.Errorf("no valid order items to process")
	}

	// 2. Inisialisasi struct Order dengan TotalAmount yang SUDAH dihitung
	order := &domain.Order{
		ID:                orderID,
		UserID:            &userID,
		CustomerName:      req.CustomerName,
		CustomerEmail:     req.CustomerEmail,
		CustomerPhone:     req.CustomerPhone,
		ShippingAddress:   &shippingAddress,
		ShippingMethod:    shippingMethod,
		ShippingCourier:   &req.ShippingCourier,
		TrackingNumber:    nil,
		TotalAmount:       totalAmount, // <-- Nilai total presisi hasil kalkulasi
		PaymentStatus:     domain.PaymentStatusPending,
		PaymentMethod:     nil,
		ExternalID:        &orderNumber,
		StockReduced:      false,
		FulfillmentStatus: domain.FulfillmentStatusNotApplicable,
	}

	// 3. Simpan Parent Order ke DB
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order header: %w", err)
	}

	// 4. Simpan Child OrderItems ke DB
	for _, item := range orderItems {
		if err := s.orderRepo.CreateItem(ctx, item); err != nil {
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}
	}

	return &domain.OrderResponse{
		ID:                order.ID,
		ExternalID:        order.ExternalID,
		UserID:            &userID,
		CustomerName:      order.CustomerName,
		CustomerEmail:     order.CustomerEmail,
		CustomerPhone:     order.CustomerPhone,
		ShippingAddress:   order.ShippingAddress,
		ShippingMethod:    string(shippingMethod),
		ShippingCourier:   order.ShippingCourier,
		TotalAmount:       totalAmount,
		PaymentStatus:     string(domain.PaymentStatusPending),
		FulfillmentStatus: string(domain.FulfillmentStatusNotApplicable),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}, nil
}

func (s *orderService) GetOrders(ctx context.Context, userID string) ([]domain.OrderResponse, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	orders, err := s.orderRepo.FindByUserID(ctx, parsedUserID)
	if err != nil {
		return nil, err
	}

	res := make([]domain.OrderResponse, 0, len(orders))
	for _, order := range orders {
		res = append(res, domain.OrderResponse{
			ID: order.ID,
			// OrderNumber:       nullableString(order.ExternalID),
			UserID:            order.UserID,
			ShippingMethod:    string(order.ShippingMethod),
			TotalAmount:       order.TotalAmount,
			PaymentStatus:     string(order.PaymentStatus),
			FulfillmentStatus: string(order.FulfillmentStatus),
			// TrackingNumber:    string(order.TrackingNumber),
			CreatedAt: order.CreatedAt,
			UpdatedAt: order.UpdatedAt,
		})
	}
	return res, nil
}

func (s *orderService) GetOrderByID(ctx context.Context, id string) (*domain.OrderResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	order, err := s.orderRepo.FindByID(ctx, parsedID)
	if err != nil {
		return nil, err
	}

	return &domain.OrderResponse{
		ID: order.ID,
		// OrderNumber:    nullableString(order.ExternalID),
		UserID:            order.UserID,
		ShippingMethod:    string(order.ShippingMethod),
		TotalAmount:       order.TotalAmount,
		PaymentStatus:     string(order.PaymentStatus),
		FulfillmentStatus: string(order.FulfillmentStatus),
		// TrackingNumber:    string(order.TrackingNumber),
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}, nil
}

func (s *orderService) UpdateShipment(ctx context.Context, orderID string, req domain.OrderShipmentRequest) (*domain.OrderResponse, error) {
	parsedID, err := uuid.Parse(orderID)
	if err != nil {
		return nil, err
	}

	order, err := s.orderRepo.FindByID(ctx, parsedID)
	if err != nil {
		return nil, err
	}

	order.TrackingNumber = &req.TrackingNumber
	order.FulfillmentStatus = domain.FulfillmentStatusShipped
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	return &domain.OrderResponse{
		ID: order.ID,
		// OrderNumber:    nullableString(order.ExternalID),
		// TrackingNumber:    req.TrackingNumber,
		FulfillmentStatus: string(order.FulfillmentStatus),
		TotalAmount:       order.TotalAmount,
		PaymentStatus:     string(order.PaymentStatus),
		CreatedAt:         order.CreatedAt,
		UpdatedAt:         order.UpdatedAt,
	}, nil
}

func (s *orderService) ProcessPaymentWebhook(ctx context.Context, payload map[string]interface{}) error {
	_ = ctx
	if payload == nil {
		return fmt.Errorf("payload is required")
	}
	return nil
}
