package domain

import (
	"time"

	"github.com/google/uuid"
)

type CartItem struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	VariantID uuid.UUID      `json:"variant_id" gorm:"type:uuid;not null"`
	Variant   ProductVariant `json:"variant,omitempty" gorm:"foreignKey:VariantID"`
	Quantity  int            `json:"quantity" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type Wishlist struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	Product   Product   `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	CreatedAt time.Time `json:"created_at"`
}

type Review struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	Rating    int       `json:"rating" gorm:"not null"`
	Comment   *string   `json:"comment,omitempty" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
}
