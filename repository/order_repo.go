package repository

import (
	"context"
	"ecommerce-catalog-api/domain"
	"fmt"
	"sort"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateWithStock(ctx context.Context, order *domain.Order, items []*domain.OrderItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		quantities := make(map[uuid.UUID]int, len(items))
		for _, item := range items {
			quantities[item.ProductVariantID] += item.Quantity
		}

		variantIDs := make([]uuid.UUID, 0, len(quantities))
		for variantID := range quantities {
			variantIDs = append(variantIDs, variantID)
		}
		sort.Slice(variantIDs, func(i, j int) bool {
			return variantIDs[i].String() < variantIDs[j].String()
		})

		for _, variantID := range variantIDs {
			quantity := quantities[variantID]
			result := tx.Model(&domain.ProductVariant{}).
				Where("id = ? AND is_active = ? AND stock >= ?", variantID, true, quantity).
				UpdateColumn("stock", gorm.Expr("stock - ?", quantity))
			if result.Error != nil {
				return fmt.Errorf("failed to reduce stock for product variant %s: %w", variantID, result.Error)
			}
			if result.RowsAffected != 1 {
				return fmt.Errorf("insufficient stock or inactive product variant: %s", variantID)
			}
		}

		order.StockReduced = true
		if err := tx.Create(order).Error; err != nil {
			return fmt.Errorf("failed to create order header: %w", err)
		}
		for _, item := range items {
			if err := tx.Create(item).Error; err != nil {
				return fmt.Errorf("failed to create order item: %w", err)
			}
		}
		return nil
	})
}

func (r *orderRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	var order domain.Order
	if err := r.db.WithContext(ctx).Preload("Items").First(&order, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Order, error) {
	var orders []domain.Order
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Preload("Items").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) Update(ctx context.Context, order *domain.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}
