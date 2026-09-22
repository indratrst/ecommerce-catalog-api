package domain

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name       string   `json:"name"`
	Price      int      `json:"price"`
	Stock      int      `json:"stock"`
	CategoryID uint     `json:"category_id"` // Foreign Key
	Category   Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

type ProductRequest struct {
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Stock      int    `json:"stock"`
	CategoryID uint   `json:"category_id"`
}

// ProductQueryParam untuk menangkap query string dari URL
type ProductQueryParam struct {
	Page       int    `query:"page"`
	Limit      int    `query:"limit"`
	Search     string `query:"search"`
	CategoryID uint   `query:"category_id"`
}

// MetaPagination metadata pagination untuk response JSON
type MetaPagination struct {
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
	Limit       int   `json:"limit"`
	TotalData   int64 `json:"total_data"`
}

// ProductPaginationResponse format response akhir
type ProductPaginationResponse struct {
	Data []Product      `json:"data"`
	Meta MetaPagination `json:"meta"`
}

// Interface Repository (Database Contract)
type ProductRepository interface {
	Create(product *Product) error
	FindAll(param ProductQueryParam) ([]Product, int64, error)
	FindByID(id string) (Product, error)
	Update(product *Product) error
	Delete(id string) error
}

// Interface Service (Business Logic Contract)
type ProductService interface {
	CreateProduct(req ProductRequest) (Product, error)
	GetAllProducts(param ProductQueryParam) (ProductPaginationResponse, error)
	GetProductByID(id string) (Product, error)
	UpdateProduct(id string, req ProductRequest) (Product, error)
	DeleteProduct(id string) error
}
