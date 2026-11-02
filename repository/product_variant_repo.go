package repository

import (
	"context"
	"ecommerce-catalog-api/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productVariantRepository struct {
	db *gorm.DB
}

func NewProductVariantRepository(db *gorm.DB) domain.ProductVariantRepository {
	return &productVariantRepository{db: db}
}

func (r *productVariantRepository) Create(ctx context.Context, variant *domain.ProductVariant) error {
	return r.db.WithContext(ctx).Create(variant).Error
}

func (r *productVariantRepository) FindByProductID(ctx context.Context, productID uuid.UUID) ([]domain.ProductVariant, error) {
	var variants []domain.ProductVariant
	err := r.db.WithContext(ctx).Where("product_id = ?", productID).Order("created_at DESC").Find(&variants).Error
	return variants, err
}

func (r *productVariantRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.ProductVariant, error) {
	var variant domain.ProductVariant
	err := r.db.WithContext(ctx).First(&variant, "id = ?", id).Error
	return &variant, err
}

func (r *productVariantRepository) Update(ctx context.Context, variant *domain.ProductVariant) error {
	return r.db.WithContext(ctx).Save(variant).Error
}

func (r *productVariantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.ProductVariant{}, "id = ?", id).Error
}
