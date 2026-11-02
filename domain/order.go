package domain

import (
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
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusSettled  PaymentStatus = "SETTLEMENT"
	PaymentStatusExpired   PaymentStatus = "EXPIRED"
	PaymentStatusCancel    PaymentStatus = "CANCEL"
	PaymentStatusFailed    PaymentStatus = "FAILED"
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
	ID               uuid.UUID        `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID           *uuid.UUID       `json:"userId,omitempty" gorm:"type:uuid;default:null"`
	User             *User            `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
	CustomerName     string           `json:"customerName" gorm:"type:varchar(255);not null"`
	CustomerEmail    string           `json:"customerEmail" gorm:"type:varchar(255);not null"`
	CustomerPhone    string           `json:"customerPhone" gorm:"type:varchar(50);not null"`
	ShippingAddress  *string          `json:"shippingAddress,omitempty" gorm:"type:text"`
	ShippingMethod   ShippingMethod   `json:"shippingMethod" gorm:"type:shipping_method_enum;not null"`
	ShippingCourier  *string          `json:"shippingCourier,omitempty" gorm:"type:varchar(100)"`
	TrackingNumber   *string          `json:"trackingNumber,omitempty" gorm:"type:varchar(100)"`
	TotalAmount      int              `json:"totalAmount" gorm:"not null"`
	PaymentStatus    PaymentStatus    `json:"paymentStatus" gorm:"type:payment_status_enum;not null;default:'PENDING'"`
	PaymentMethod    *string          `json:"paymentMethod,omitempty" gorm:"type:varchar(100)"`
	SnapToken        *string          `json:"snapToken,omitempty" gorm:"type:varchar(255)"`
	SnapRedirectUrl  *string          `json:"snapRedirectUrl,omitempty" gorm:"type:text"`
	ExternalID       *string          `json:"externalId,omitempty" gorm:"type:varchar(255);unique"`
	StockReduced     bool             `json:"stockReduced" gorm:"default:false"`
	FulfillmentStatus FulfillmentStatus `json:"fulfillmentStatus" gorm:"type:fulfillment_status_enum;not null;default:'NOT_APPLICABLE'"`
	CreatedAt        time.Time        `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt        time.Time        `json:"updatedAt" gorm:"autoUpdateTime"`
	Items            []OrderItem     `json:"items,omitempty" gorm:"foreignKey:OrderID;references:ID"`
}

type OrderItem struct {
	ID               uuid.UUID              `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrderID          uuid.UUID              `json:"orderId" gorm:"type:uuid;not null"`
	Order            *Order                 `json:"order,omitempty" gorm:"foreignKey:OrderID;references:ID"`
	ProductVariantID uuid.UUID              `json:"productVariantId" gorm:"type:uuid;not null"`
	Quantity         int                    `json:"quantity" gorm:"not null"`
	PriceAtPurchase  int                    `json:"priceAtPurchase" gorm:"not null"`
	VariantSnapshot  map[string]interface{} `json:"variantSnapshot,omitempty" gorm:"type:jsonb"`
}
