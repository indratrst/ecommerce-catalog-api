package domain

import "time"

type CouponType string

const (
	CouponPercentage  CouponType = "percentage"
	CouponFixedAmount CouponType = "fixed_amount"
)

type Coupon struct {
	ID            uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Code          string     `json:"code" gorm:"type:varchar(50);not null;unique"`
	DiscountType  CouponType `json:"discount_type" gorm:"type:coupon_type;not null"`
	DiscountValue float64    `json:"discount_value" gorm:"type:numeric(12,2);not null"`
	MinPurchase   float64    `json:"min_purchase" gorm:"type:numeric(12,2);default:0.0"`
	MaxDiscount   *float64   `json:"max_discount,omitempty" gorm:"type:numeric(12,2)"`
	ExpiredAt     time.Time  `json:"expired_at" gorm:"not null"`
	IsActive      bool       `json:"is_active" gorm:"default:true"`
	CreatedAt     time.Time  `json:"created_at"`
}
