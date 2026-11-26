package service

import (
	"context"
	"ecommerce-catalog-api/domain"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
)

type orderService struct {
	orderRepo domain.OrderRepository
	midtrans  *midtransClient
}

// nullableString mengubah pointer string menjadi string kosong jika nil.
func nullableString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

// NewOrderService membuat layanan order dengan repository yang digunakan untuk
// menyimpan dan mengambil data order.
func NewOrderService(orderRepo domain.OrderRepository) domain.OrderService {
	return &orderService{orderRepo: orderRepo, midtrans: newMidtransClientFromEnv()}
}

// CreateOrder memvalidasi data checkout, menghitung total dari item, lalu
// menyimpan order dan item-itemnya. UserID bersifat opsional untuk checkout
// tamu; alamat pengiriman dapat menggunakan StoreLocation jika tidak diisi.
func (s *orderService) CreateOrder(ctx context.Context, req domain.OrderRequest) (*domain.OrderResponse, error) {
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("items are required")
	}

	req.CustomerName = strings.TrimSpace(req.CustomerName)
	req.CustomerEmail = strings.TrimSpace(req.CustomerEmail)
	req.CustomerPhone = strings.TrimSpace(req.CustomerPhone)
	if req.CustomerName == "" {
		return nil, fmt.Errorf("customerName is required")
	}
	if req.CustomerEmail == "" {
		return nil, fmt.Errorf("customerEmail is required")
	}
	if _, err := mail.ParseAddress(req.CustomerEmail); err != nil {
		return nil, fmt.Errorf("invalid customerEmail")
	}
	if req.CustomerPhone == "" {
		return nil, fmt.Errorf("customerPhone is required")
	}

	var userID *uuid.UUID
	if strings.TrimSpace(req.UserID) != "" {
		parsedUserID, err := uuid.Parse(req.UserID)
		if err != nil || parsedUserID == uuid.Nil {
			return nil, fmt.Errorf("invalid userId")
		}
		userID = &parsedUserID
	}

	deliveryType := strings.ToLower(strings.TrimSpace(req.DeliveryType))
	requestedShippingMethod := strings.ToUpper(strings.TrimSpace(req.ShippingMethod))
	var shippingMethod domain.ShippingMethod
	if deliveryType == "pickup" || requestedShippingMethod == string(domain.ShippingMethodPickupStore) {
		shippingMethod = domain.ShippingMethodPickupStore
	} else if deliveryType == "courier" || requestedShippingMethod == string(domain.ShippingMethodShipping) {
		shippingMethod = domain.ShippingMethodShipping
	} else {
		return nil, fmt.Errorf("valid shippingMethod is required")
	}

	orderID := uuid.New()
	orderNumber := fmt.Sprintf("INV-%s-%s", time.Now().Format("20060102"), orderID)
	shippingAddress := strings.TrimSpace(req.ShippingAddress)
	if shippingAddress == "" {
		shippingAddress = strings.TrimSpace(req.StoreLocation)
	}
	if shippingAddress == "" {
		return nil, fmt.Errorf("shippingAddress or storeLocation is required")
	}
	shippingCourier := strings.TrimSpace(req.ShippingCourier)
	if shippingMethod == domain.ShippingMethodShipping && shippingCourier == "" {
		return nil, fmt.Errorf("shippingCourier is required for shipping orders")
	}
	if req.ShippingCost < 0 {
		return nil, fmt.Errorf("shippingCost cannot be negative")
	}
	if shippingMethod == domain.ShippingMethodPickupStore && req.ShippingCost != 0 {
		return nil, fmt.Errorf("pickup orders cannot have a shippingCost")
	}
	var shippingCourierValue *string
	if shippingCourier != "" {
		shippingCourierValue = &shippingCourier
	}

	var orderItems []*domain.OrderItem
	var totalAmount int = 0

	// 1. Loop items & hitung totalAmount TERLEBIH DAHULU
	for _, itemReq := range req.Items {
		if itemReq.Quantity <= 0 {
			return nil, fmt.Errorf("item quantity must be greater than zero")
		}

		// Fix: Ambil ProductVariantID, jika kosong baru fallback ke VariantID
		variantIDStr := strings.TrimSpace(itemReq.ProductVariantID)
		if variantIDStr == "" {
			variantIDStr = strings.TrimSpace(itemReq.VariantID)
		}

		// Jika keduanya tetap kosong, lempar error jelas
		if variantIDStr == "" {
			return nil, fmt.Errorf("productVariantId / variantId is required")
		}

		variantID, err := uuid.Parse(variantIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid productVariantId: %s", variantIDStr)
		}
		if itemReq.PriceAtPurchase < 0 {
			return nil, fmt.Errorf("priceAtPurchase cannot be negative")
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
	totalAmount += req.ShippingCost

	// 2. Inisialisasi struct Order dengan TotalAmount yang SUDAH dihitung
	order := &domain.Order{
		ID:                orderID,
		UserID:            userID,
		CustomerName:      req.CustomerName,
		CustomerEmail:     req.CustomerEmail,
		CustomerPhone:     req.CustomerPhone,
		ShippingAddress:   &shippingAddress,
		ShippingMethod:    shippingMethod,
		ShippingCourier:   shippingCourierValue,
		TrackingNumber:    nil,
		TotalAmount:       totalAmount, // <-- Nilai total presisi hasil kalkulasi
		PaymentStatus:     domain.PaymentStatusPending,
		PaymentMethod:     nil,
		ExternalID:        &orderNumber,
		StockReduced:      false,
		FulfillmentStatus: domain.FulfillmentStatusNotApplicable,
	}

	if strings.TrimSpace(s.midtrans.serverKey) == "" {
		return nil, fmt.Errorf("MIDTRANS_SERVER_KEY is required")
	}

	// Simpan order tanpa memotong stok; stok baru berubah setelah pembayaran berhasil.
	if err := s.orderRepo.CreateWithItems(ctx, order, orderItems); err != nil {
		return nil, err
	}

	transaction, err := s.midtrans.createTransaction(ctx, order, orderItems, req.ShippingCost)
	if err != nil {
		return nil, err
	}
	order.SnapToken = &transaction.Token
	order.SnapRedirectUrl = &transaction.RedirectURL
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("Midtrans transaction created but failed to save payment details: %w", err)
	}

	return &domain.OrderResponse{
		ID:                order.ID,
		ExternalID:        order.ExternalID,
		UserID:            userID,
		CustomerName:      order.CustomerName,
		CustomerEmail:     order.CustomerEmail,
		CustomerPhone:     order.CustomerPhone,
		ShippingAddress:   order.ShippingAddress,
		ShippingMethod:    string(shippingMethod),
		ShippingCourier:   order.ShippingCourier,
		TotalAmount:       totalAmount,
		PaymentStatus:     string(domain.PaymentStatusPending),
		PaymentMethod:     order.PaymentMethod,
		StockReduced:      order.StockReduced,
		FulfillmentStatus: string(domain.FulfillmentStatusNotApplicable),
		SnapToken:         order.SnapToken,
		SnapRedirectUrl:   order.SnapRedirectUrl,
		CreatedAt:         order.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:         order.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// GetOrders mengambil seluruh order milik pengguna berdasarkan UUID userID,
// lalu mengubahnya menjadi daftar OrderResponse.
func (s *orderService) GetOrders(ctx context.Context, userID string) ([]domain.OrderResponse, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	orders, err := s.orderRepo.FindByUserID(ctx, parsedUserID)
	if err != nil {
		return nil, err
	}

	return orderResponses(orders), nil
}

func (s *orderService) GetAllOrders(ctx context.Context) ([]domain.OrderResponse, error) {
	orders, err := s.orderRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return orderResponses(orders), nil
}

// GetOrderByID mengambil order berdasarkan UUID id dan mengubahnya menjadi
// OrderResponse.
func (s *orderService) GetOrderByID(ctx context.Context, id string) (*domain.OrderResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	order, err := s.orderRepo.FindByID(ctx, parsedID)
	if err != nil {
		return nil, err
	}

	response := orderResponses([]domain.Order{*order})[0]
	return &response, nil
}

// UpdateShipment menyimpan nomor pelacakan pada order dan mengubah status
// pemenuhannya menjadi shipped. orderID harus berupa UUID yang valid.
func (s *orderService) UpdateShipment(ctx context.Context, orderID string, req domain.OrderShipmentRequest) (*domain.OrderResponse, error) {
	parsedID, err := uuid.Parse(orderID)
	if err != nil {
		return nil, err
	}

	order, err := s.orderRepo.FindByID(ctx, parsedID)
	if err != nil {
		return nil, err
	}

	switch domain.FulfillmentStatus(req.Status) {
	case domain.FulfillmentStatusReadyToPickup:
		if order.ShippingMethod != domain.ShippingMethodPickupStore {
			return nil, fmt.Errorf("ready to pickup only applies to pickup orders")
		}
		order.FulfillmentStatus = domain.FulfillmentStatusReadyToPickup
	case domain.FulfillmentStatusPickedUp:
		if order.ShippingMethod != domain.ShippingMethodPickupStore {
			return nil, fmt.Errorf("picked up only applies to pickup orders")
		}
		order.FulfillmentStatus = domain.FulfillmentStatusPickedUp
	case domain.FulfillmentStatusShipped:
		if strings.TrimSpace(req.TrackingNumber) == "" {
			return nil, fmt.Errorf("tracking number is required")
		}
		trackingNumber := strings.TrimSpace(req.TrackingNumber)
		order.TrackingNumber = &trackingNumber
		order.FulfillmentStatus = domain.FulfillmentStatusShipped
	case domain.FulfillmentStatusDelivered:
		order.FulfillmentStatus = domain.FulfillmentStatusDelivered
	default:
		return nil, fmt.Errorf("unsupported fulfillment status %q", req.Status)
	}
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	response := orderResponses([]domain.Order{*order})[0]
	return &response, nil
}

func orderResponses(orders []domain.Order) []domain.OrderResponse {
	responses := make([]domain.OrderResponse, 0, len(orders))
	for _, order := range orders {
		responses = append(responses, domain.OrderResponse{
			ID: order.ID, ExternalID: order.ExternalID, UserID: order.UserID,
			CustomerName: order.CustomerName, CustomerEmail: order.CustomerEmail,
			CustomerPhone: order.CustomerPhone, ShippingAddress: order.ShippingAddress,
			ShippingMethod: string(order.ShippingMethod), ShippingCourier: order.ShippingCourier,
			TrackingNumber: order.TrackingNumber, TotalAmount: order.TotalAmount,
			PaymentStatus: string(order.PaymentStatus), PaymentMethod: order.PaymentMethod,
			StockReduced: order.StockReduced, FulfillmentStatus: string(order.FulfillmentStatus),
			SnapToken: order.SnapToken, SnapRedirectUrl: order.SnapRedirectUrl,
			CreatedAt: order.CreatedAt.Format(time.RFC3339),
			UpdatedAt: order.UpdatedAt.Format(time.RFC3339), Items: order.Items,
		})
	}
	return responses
}

func (s *orderService) ProcessPaymentWebhook(ctx context.Context, payload map[string]interface{}) error {
	if payload == nil {
		return fmt.Errorf("payload is required")
	}
	if orderID, ok := payload["order_id"].(string); ok && strings.HasPrefix(orderID, midtransTestOrderIDPrefix) {
		return nil
	}
	externalID, amount, paymentMethod, status, err := s.midtrans.verifyNotification(payload)
	if err != nil {
		return err
	}
	if status == domain.PaymentStatusPending {
		return nil
	}
	return s.orderRepo.ProcessPaymentNotification(ctx, externalID, amount, status, paymentMethod)
}
