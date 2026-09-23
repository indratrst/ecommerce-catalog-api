package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Entity Product
type Product struct {
	ID          uuid.UUID        `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CategoryID  uuid.UUID        `json:"category_id" gorm:"type:uuid;not null"`
	Category    Category         `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Name        string           `json:"name" gorm:"type:varchar(150);not null"`
	Slug        string           `json:"slug" gorm:"type:varchar(150);not null;unique"`
	Description string           `json:"description" gorm:"type:text"`
	Price       float64          `json:"price" gorm:"type:decimal(12,2);not null"`
	Stock       int              `json:"stock" gorm:"type:int;default:0"`
	Variants    []ProductVariant `json:"variants,omitempty" gorm:"foreignKey:ProductID"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// DTO Request
type ProductRequest struct {
	CategoryID  string  `json:"category_id" validate:"required,uuid"`
	Name        string  `json:"name" validate:"required,min=3,max=150"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"gte=0"`
}

// Query Parameter DTO (Searching, Filtering, Pagination)
type ProductQueryParam struct {
	Search     string `query:"search"`
	CategoryID string `query:"category_id"`
	Page       int    `query:"page"`
	Limit      int    `query:"limit"`
}

// DTO Response
type ProductResponse struct {
	ID          uuid.UUID         `json:"id"`
	CategoryID  uuid.UUID         `json:"category_id"`
	Category    *CategoryResponse `json:"category,omitempty"`
	Name        string            `json:"name"`
	Slug        string            `json:"slug"`
	Description string            `json:"description"`
	Price       float64           `json:"price"`
	Stock       int               `json:"stock"`
	CreatedAt   time.Time         `json:"created_at"`
}

// Pagination Response
type ProductPaginationResponse struct {
	Data        []ProductResponse `json:"data"`
	TotalData   int64             `json:"total_data"`
	TotalPages  int               `json:"total_pages"`
	CurrentPage int               `json:"current_page"`
	Limit       int               `json:"limit"`
}

// Contract Repository
type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	FindAll(ctx context.Context, param ProductQueryParam) ([]Product, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// Contract Service
type ProductService interface {
	CreateProduct(ctx context.Context, req ProductRequest) (*ProductResponse, error)
	GetProducts(ctx context.Context, param ProductQueryParam) (*ProductPaginationResponse, error)
	GetProductByID(ctx context.Context, id string) (*ProductResponse, error)
	UpdateProduct(ctx context.Context, id string, req ProductRequest) (*ProductResponse, error)
	DeleteProduct(ctx context.Context, id string) error
}

// Entity Product Variant
type ProductVariant struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	Product   Product   `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	SKU       string    `json:"sku" gorm:"type:varchar(100);not null;unique"`
	Name      string    `json:"name" gorm:"type:varchar(100);not null"` // Contoh: "S / Red" atau "XL"
	Price     float64   `json:"price" gorm:"type:decimal(12,2);not null"`
	Stock     int       `json:"stock" gorm:"type:int;default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DTO Response untuk ProductVariant
type ProductVariantResponse struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	SKU       string    `json:"sku"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
}

// Contract Repository untuk ProductVariant
type ProductVariantRepository interface {
	Create(ctx context.Context, variant *ProductVariant) error
	FindByProductID(ctx context.Context, productID uuid.UUID) ([]ProductVariant, error)
	FindByID(ctx context.Context, id uuid.UUID) (*ProductVariant, error)
	Update(ctx context.Context, variant *ProductVariant) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type Size struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"type:varchar(50);not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Color struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"type:varchar(50);not null"`
	HexCode   string    `json:"hex_code" gorm:"type:varchar(10)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductImage struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	ImageURL  string    `json:"image_url" gorm:"type:text;not null"`
	IsPrimary bool      `json:"is_primary" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
}
