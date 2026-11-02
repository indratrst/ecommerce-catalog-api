package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Entity Category
type Category struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null"`
	Slug        string    `json:"slug" gorm:"type:varchar(255);not null;unique"`
	Description *string   `json:"description,omitempty" gorm:"type:text"`
	Image       *string   `json:"image,omitempty" gorm:"type:text"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

// DTO Request
type CategoryRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=255"`
	Description *string `json:"description,omitempty"`
	Image       *string `json:"image,omitempty"`
}

// DTO Response
type CategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description *string   `json:"description,omitempty"`
	Image       *string   `json:"image,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Contract Repository
type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	FindAll(ctx context.Context) ([]Category, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// Contract Service
type CategoryService interface {
	CreateCategory(ctx context.Context, req CategoryRequest) (*CategoryResponse, error)
	GetCategories(ctx context.Context) ([]CategoryResponse, error)
	GetCategoryByID(ctx context.Context, id string) (*CategoryResponse, error)
}
