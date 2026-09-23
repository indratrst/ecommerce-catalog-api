package repository

import (
	"context"
	"ecommerce-catalog-api/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) FindAll(ctx context.Context, param domain.ProductQueryParam) ([]domain.Product, int64, error) {
	var products []domain.Product
	var totalData int64

	query := r.db.WithContext(ctx).Model(&domain.Product{}).Preload("Category")

	if param.Search != "" {
		query = query.Where("name ILIKE ?", "%"+param.Search+"%")
	}

	if param.CategoryID != "" {
		query = query.Where("category_id = ?", param.CategoryID)
	}

	if err := query.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	page := param.Page
	if page <= 0 {
		page = 1
	}
	limit := param.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	err := query.Limit(limit).Offset(offset).Find(&products).Error
	return products, totalData, err
}
func (r *productRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	var product domain.Product
	err := r.db.WithContext(ctx).Preload("Category").First(&product, "id = ?", id).Error
	return &product, err
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Product{}, "id = ?", id).Error
}
