package domain

import (
	"time"

	"github.com/google/uuid"
)

type DeliveryType string

const (
	DeliveryCourier DeliveryType = "courier"
	DeliveryPickup  DeliveryType = "pickup"
)

type OrderStatus string

const (
	StatusPendingPayment OrderStatus = "pending_payment"
	StatusPaid           OrderStatus = "paid"
	StatusProcessing     OrderStatus = "processing"
	StatusReadyForPickup OrderStatus = "ready_for_pickup"
	StatusShipped        OrderStatus = "shipped"
	StatusCompleted      OrderStatus = "completed"
	StatusCancelled      OrderStatus = "cancelled"
)

type Order struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrderNumber string    `json:"order_number" gorm:"type:varchar(50);not null;unique"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	User        User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	CouponID    *uint     `json:"coupon_id,omitempty"`
	Coupon      *Coupon   `json:"coupon,omitempty" gorm:"foreignKey:CouponID"`

	// Mode Pengiriman (Courier vs Pickup in Store)
	DeliveryType  DeliveryType `json:"delivery_type" gorm:"type:delivery_type;not null;default:'courier'"`
	StoreLocation *string      `json:"store_location,omitempty" gorm:"type:varchar(150)"`
	PickupCode    *string      `json:"pickup_code,omitempty" gorm:"type:varchar(20)"`

	// Informasi Kurir (Opsional jika pickup)
	ShippingRecipient *string `json:"shipping_recipient,omitempty" gorm:"type:varchar(100)"`
	ShippingPhone     *string `json:"shipping_phone,omitempty" gorm:"type:varchar(20)"`
	ShippingAddress   *string `json:"shipping_address,omitempty" gorm:"type:text"`
	ShippingCourier   *string `json:"shipping_courier,omitempty" gorm:"type:varchar(50)"`
	ShippingService   *string `json:"shipping_service,omitempty" gorm:"type:varchar(50)"`
	TrackingNumber    *string `json:"tracking_number,omitempty" gorm:"type:varchar(100)"`

	// Financial Totals
	Subtotal       float64 `json:"subtotal" gorm:"type:numeric(12,2);not null"`
	DiscountAmount float64 `json:"discount_amount" gorm:"type:numeric(12,2);default:0.0"`
	ShippingCost   float64 `json:"shipping_cost" gorm:"type:numeric(12,2);default:0.0"`
	TotalAmount    float64 `json:"total_amount" gorm:"type:numeric(12,2);not null"`

	// Payment Gateway Integrations
	PaymentToken *string     `json:"payment_token,omitempty" gorm:"type:varchar(255)"`
	PaymentURL   *string     `json:"payment_url,omitempty" gorm:"type:text"`
	Status       OrderStatus `json:"status" gorm:"type:order_status;default:'pending_payment'"`

	Items     []OrderItem `json:"items,omitempty" gorm:"foreignKey:OrderID"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID           uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID      uuid.UUID      `json:"order_id" gorm:"type:uuid;not null"`
	VariantID    uuid.UUID      `json:"variant_id" gorm:"type:uuid;not null"`
	Variant      ProductVariant `json:"variant,omitempty" gorm:"foreignKey:VariantID"`
	ProductName  string         `json:"product_name" gorm:"type:varchar(255);not null"`
	VariantSize  string         `json:"variant_size" gorm:"type:varchar(20);not null"`
	VariantColor *string        `json:"variant_color,omitempty" gorm:"type:varchar(50)"`
	Price        float64        `json:"price" gorm:"type:numeric(12,2);not null"`
	Quantity     int            `json:"quantity" gorm:"not null"`
	Subtotal     float64        `json:"subtotal" gorm:"type:numeric(12,2);not null"`
}
