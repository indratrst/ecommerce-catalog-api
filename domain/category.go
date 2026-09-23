package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Entity Category
type Category struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ParentID  *uint      `json:"parent_id,omitempty"`
	Name      string     `json:"name" gorm:"type:varchar(100);not null"`
	Slug      string     `json:"slug" gorm:"type:varchar(100);not null;unique"`
	ImageURL  *string    `json:"image_url,omitempty" gorm:"type:text"`
	Children  []Category `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	CreatedAt time.Time  `json:"created_at"`
}

// DTO Request
type CategoryRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

// DTO Response
type CategoryResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
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
