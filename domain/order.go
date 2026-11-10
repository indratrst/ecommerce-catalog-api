package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ShippingMethod string

const (
	ShippingMethodPickupStore ShippingMethod = "PICKUP_STORE"
	ShippingMethodShipping    ShippingMethod = "SHIPPING"
)

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "PENDING"
	PaymentStatusSettled PaymentStatus = "SETTLEMENT"
	PaymentStatusExpired PaymentStatus = "EXPIRED"
	PaymentStatusCancel  PaymentStatus = "CANCEL"
	PaymentStatusFailed  PaymentStatus = "FAILED"
)

type FulfillmentStatus string

const (
	FulfillmentStatusNotApplicable FulfillmentStatus = "NOT_APPLICABLE"
	FulfillmentStatusProcessing    FulfillmentStatus = "PROCESSING"
	FulfillmentStatusReadyToPickup FulfillmentStatus = "READY_TO_PICKUP"
	FulfillmentStatusPickedUp      FulfillmentStatus = "PICKED_UP"
	FulfillmentStatusShipped       FulfillmentStatus = "SHIPPED"
	FulfillmentStatusDelivered     FulfillmentStatus = "DELIVERED"
)

type Order struct {
	ID                uuid.UUID         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID            *uuid.UUID        `json:"userId,omitempty" gorm:"type:uuid;default:null"`
	User              *User             `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
	CustomerName      string            `json:"customerName" gorm:"type:varchar(255);not null"`
	CustomerEmail     string            `json:"customerEmail" gorm:"type:varchar(255);not null"`
	CustomerPhone     string            `json:"customerPhone" gorm:"type:varchar(50);not null"`
	ShippingAddress   *string           `json:"shippingAddress,omitempty" gorm:"type:text"`
	ShippingMethod    ShippingMethod    `json:"shippingMethod" gorm:"type:shipping_method_enum;not null"`
	ShippingCourier   *string           `json:"shippingCourier,omitempty" gorm:"type:varchar(100)"`
	TrackingNumber    *string           `json:"trackingNumber,omitempty" gorm:"type:varchar(100)"`
	TotalAmount       int               `json:"totalAmount" gorm:"not null"`
	PaymentStatus     PaymentStatus     `json:"paymentStatus" gorm:"type:payment_status_enum;not null;default:'PENDING'"`
	PaymentMethod     *string           `json:"paymentMethod,omitempty" gorm:"type:varchar(100)"`
	SnapToken         *string           `json:"snapToken,omitempty" gorm:"type:varchar(255)"`
	SnapRedirectUrl   *string           `json:"snapRedirectUrl,omitempty" gorm:"type:text"`
	ExternalID        *string           `json:"externalId,omitempty" gorm:"type:varchar(255);unique"`
	StockReduced      bool              `json:"stockReduced" gorm:"default:false"`
	FulfillmentStatus FulfillmentStatus `json:"fulfillmentStatus" gorm:"type:fulfillment_status_enum;not null;default:'NOT_APPLICABLE'"`
	CreatedAt         time.Time         `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt         time.Time         `json:"updatedAt" gorm:"autoUpdateTime"`
	Items             []OrderItem       `json:"items,omitempty" gorm:"foreignKey:OrderID;references:ID;constraint:OnDelete:CASCADE"`
}

type OrderItem struct {
	ID               uuid.UUID              `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrderID          uuid.UUID              `json:"orderId" gorm:"type:uuid;not null"`
	Order            *Order                 `json:"order,omitempty" gorm:"foreignKey:OrderID;references:ID"`
	ProductVariantID uuid.UUID              `json:"productVariantId" gorm:"type:uuid;not null"`
	ProductVariant   *ProductVariant        `json:"productVariant,omitempty" gorm:"foreignKey:ProductVariantID;references:ID"`
	Quantity         int                    `json:"quantity" gorm:"not null"`
	PriceAtPurchase  int                    `json:"priceAtPurchase" gorm:"not null"`
	VariantSnapshot  map[string]interface{} `json:"variantSnapshot,omitempty" gorm:"type:jsonb;serializer:json"`
}

// Order request/response DTOs

type OrderItemRequest struct {
	ProductVariantID string `json:"productVariantId"`
	VariantID        string `json:"variantId,omitempty"` // Alias untuk kompatibilitas
	Quantity         int    `json:"quantity"`
	PriceAtPurchase  int    `json:"priceAtPurchase,omitempty"`
}

type OrderRequest struct {
	UserID          string             `json:"userId,omitempty"`
	CustomerName    string             `json:"customerName"`
	CustomerEmail   string             `json:"customerEmail"`
	CustomerPhone   string             `json:"customerPhone"`
	DeliveryType    string             `json:"deliveryType,omitempty"`
	ShippingMethod  string             `json:"shippingMethod"`
	ShippingCourier string             `json:"shippingCourier,omitempty"`
	ShippingAddress string             `json:"shippingAddress,omitempty"`
	StoreLocation   string             `json:"storeLocation,omitempty"`
	TotalAmount     int                `json:"totalAmount,omitempty"`
	Items           []OrderItemRequest `json:"items"`
}

type OrderShipmentRequest struct {
	TrackingNumber string `json:"trackingNumber"`
	Status         string `json:"status"`
}

type OrderResponse struct {
	ID                uuid.UUID   `json:"id"`
	ExternalID        *string     `json:"externalId,omitempty"`
	UserID            *uuid.UUID  `json:"userId,omitempty"`
	CustomerName      string      `json:"customerName"`
	CustomerEmail     string      `json:"customerEmail"`
	CustomerPhone     string      `json:"customerPhone"`
	ShippingAddress   *string     `json:"shippingAddress,omitempty"`
	ShippingMethod    string      `json:"shippingMethod"`
	ShippingCourier   *string     `json:"shippingCourier,omitempty"`
	TrackingNumber    *string     `json:"trackingNumber,omitempty"`
	TotalAmount       int         `json:"totalAmount"`
	PaymentStatus     string      `json:"paymentStatus"`
	FulfillmentStatus string      `json:"fulfillmentStatus"`
	SnapToken         *string     `json:"snapToken,omitempty"`
	SnapRedirectUrl   *string     `json:"snapRedirectUrl,omitempty"`
	CreatedAt         string      `json:"createdAt"`
	UpdatedAt         string      `json:"updatedAt"`
	Items             []OrderItem `json:"items,omitempty"`
}

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	CreateItem(ctx context.Context, item *OrderItem) error
	FindByID(ctx context.Context, id uuid.UUID) (*Order, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]Order, error)
	Update(ctx context.Context, order *Order) error
}

type OrderService interface {
	CreateOrder(ctx context.Context, req OrderRequest) (*OrderResponse, error)
	GetOrders(ctx context.Context, userID string) ([]OrderResponse, error)
	GetOrderByID(ctx context.Context, id string) (*OrderResponse, error)
	UpdateShipment(ctx context.Context, orderID string, req OrderShipmentRequest) (*OrderResponse, error)
	ProcessPaymentWebhook(ctx context.Context, payload map[string]interface{}) error
}
