package repository

import (
	"context"
	"ecommerce-catalog-api/domain"
	"fmt"
	"sort"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateWithItems(ctx context.Context, order *domain.Order, items []*domain.OrderItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order.StockReduced = false
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

func (r *orderRepository) ProcessPaymentNotification(ctx context.Context, externalID string, amount int, status domain.PaymentStatus, paymentMethod string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order domain.Order
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("external_id = ?", externalID).
			First(&order).Error; err != nil {
			return fmt.Errorf("failed to find order for payment notification: %w", err)
		}
		if order.TotalAmount != amount {
			return fmt.Errorf("payment amount mismatch for order %s", externalID)
		}

		if order.StockReduced {
			return nil
		}

		if status == domain.PaymentStatusSettled {
			quantities := make(map[uuid.UUID]int, len(order.Items))
			if err := tx.Where("order_id = ?", order.ID).Find(&order.Items).Error; err != nil {
				return fmt.Errorf("failed to load order items: %w", err)
			}
			for _, item := range order.Items {
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
		} else if order.PaymentStatus == domain.PaymentStatusSettled {
			return nil
		}

		updates := map[string]interface{}{
			"payment_status": status,
			"stock_reduced":  status == domain.PaymentStatusSettled,
		}
		if paymentMethod != "" {
			updates["payment_method"] = paymentMethod
		}

		if status == domain.PaymentStatusSettled &&
			(order.FulfillmentStatus == "" || order.FulfillmentStatus == domain.FulfillmentStatusNotApplicable) {
			if order.ShippingMethod == domain.ShippingMethodPickupStore {
				updates["fulfillment_status"] = domain.FulfillmentStatusReadyToPickup
			} else {
				updates["fulfillment_status"] = domain.FulfillmentStatusProcessing
			}
		}

		if err := tx.Model(&order).Updates(updates).Error; err != nil {
			return fmt.Errorf("failed to update order payment status: %w", err)
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
