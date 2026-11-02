package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Entity Product
type Product struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CategoryID    uuid.UUID `json:"categoryId" gorm:"type:uuid;not null"`
	Category      *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID;references:ID"`
	Title         string    `json:"title" gorm:"type:varchar(255);not null"`
	Price         int       `json:"price" gorm:"not null"`
	OriginalPrice *int      `json:"originalPrice,omitempty"`
	Description   string    `json:"description" gorm:"type:text;not null"`
	Image         *string   `json:"image,omitempty" gorm:"type:text"`
	RatingRate    float64   `json:"ratingRate" gorm:"default:0"`
	RatingCount   int       `json:"ratingCount" gorm:"default:0"`
	CreatedAt     time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

// DTO Request
type ProductRequest struct {
	CategoryID    string  `json:"categoryId" validate:"required,uuid"`
	Title         string  `json:"title" validate:"required,min=3,max=255"`
	Price         int     `json:"price" validate:"required,gt=0"`
	OriginalPrice *int    `json:"originalPrice,omitempty"`
	Description   string  `json:"description" validate:"required"`
	Image         *string `json:"image,omitempty"`
}

// Query Parameter DTO (Searching, Filtering, Pagination)
type ProductQueryParam struct {
	Search     string `query:"search"`
	CategoryID string `query:"categoryId"`
	Page       int    `query:"page"`
	Limit      int    `query:"limit"`
}

// DTO Response
type ProductResponse struct {
	ID            uuid.UUID         `json:"id"`
	CategoryID    uuid.UUID         `json:"categoryId"`
	Category      *CategoryResponse `json:"category,omitempty"`
	Title         string            `json:"title"`
	Price         int               `json:"price"`
	OriginalPrice *int              `json:"originalPrice,omitempty"`
	Description   string            `json:"description"`
	Image         *string           `json:"image,omitempty"`
	RatingRate    float64           `json:"ratingRate"`
	RatingCount   int               `json:"ratingCount"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
}

// Pagination Response
type ProductPaginationResponse struct {
	Data        []ProductResponse `json:"data"`
	TotalData   int64             `json:"totalData"`
	TotalPages  int               `json:"totalPages"`
	CurrentPage int               `json:"currentPage"`
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
	ProductID uuid.UUID `json:"productId" gorm:"type:uuid;not null"`
	Product   *Product  `json:"product,omitempty" gorm:"foreignKey:ProductID;references:ID"`
	Size      string    `json:"size" gorm:"type:varchar(50);not null"`
	Color     *string   `json:"color,omitempty" gorm:"type:varchar(50)"`
	Stock     int       `json:"stock" gorm:"type:int;not null;default:0"`
	IsActive  bool      `json:"isActive" gorm:"default:true"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

// DTO Response untuk ProductVariant
type ProductVariantResponse struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"productId"`
	Size      string    `json:"size"`
	Color     *string   `json:"color,omitempty"`
	Stock     int       `json:"stock"`
	IsActive  bool      `json:"isActive"`
}

// DTO Request untuk ProductVariant
type ProductVariantRequest struct {
	Size     string  `json:"size" validate:"required"`
	Color    *string `json:"color,omitempty"`
	Stock    int     `json:"stock" validate:"gte=0"`
	IsActive *bool   `json:"isActive,omitempty"`
}

// Contract Repository untuk ProductVariant
type ProductVariantRepository interface {
	Create(ctx context.Context, variant *ProductVariant) error
	FindByProductID(ctx context.Context, productID uuid.UUID) ([]ProductVariant, error)
	FindByID(ctx context.Context, id uuid.UUID) (*ProductVariant, error)
	Update(ctx context.Context, variant *ProductVariant) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// Contract Service untuk ProductVariant
type ProductVariantService interface {
	CreateProductVariant(ctx context.Context, productID string, req ProductVariantRequest) (*ProductVariantResponse, error)
	GetProductVariantsByProductID(ctx context.Context, productID string) ([]ProductVariantResponse, error)
	GetProductVariantByID(ctx context.Context, id string) (*ProductVariantResponse, error)
	UpdateProductVariant(ctx context.Context, id string, req ProductVariantRequest) (*ProductVariantResponse, error)
	DeleteProductVariant(ctx context.Context, id string) error
}
