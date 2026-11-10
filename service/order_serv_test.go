package service

import (
	"context"
	"testing"
	"time"

	"ecommerce-catalog-api/domain"

	"github.com/google/uuid"
)

type orderRepositoryStub struct {
	order *domain.Order
	items []*domain.OrderItem
}

func (r *orderRepositoryStub) Create(_ context.Context, order *domain.Order) error {
	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = now
	r.order = order
	return nil
}

func (r *orderRepositoryStub) CreateItem(_ context.Context, item *domain.OrderItem) error {
	r.items = append(r.items, item)
	return nil
}

func (r *orderRepositoryStub) FindByID(context.Context, uuid.UUID) (*domain.Order, error) {
	return nil, nil
}

func (r *orderRepositoryStub) FindByUserID(context.Context, uuid.UUID) ([]domain.Order, error) {
	return nil, nil
}

func (r *orderRepositoryStub) Update(context.Context, *domain.Order) error {
	return nil
}

func TestCreateOrderAllowsGuestCheckout(t *testing.T) {
	repo := &orderRepositoryStub{}
	orderService := NewOrderService(repo)
	req := domain.OrderRequest{
		CustomerName:    "Muhamad Nur Indra",
		CustomerEmail:   "indra@example.com",
		CustomerPhone:   "081234567890",
		ShippingAddress: "Jl. Cibabat No. 123, Cimahi",
		ShippingMethod:  "SHIPPING",
		ShippingCourier: "JNE",
		TotalAmount:     1,
		Items: []domain.OrderItemRequest{{
			ProductVariantID: "b58a1c7d-a613-42a7-ae41-df544fac92a9",
			Quantity:         1,
			PriceAtPurchase:  350000,
		}},
	}

	response, err := orderService.CreateOrder(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateOrder() error = %v", err)
	}
	if repo.order == nil {
		t.Fatal("CreateOrder() did not save an order")
	}
	if repo.order.UserID != nil || response.UserID != nil {
		t.Fatal("guest order should not be associated with a user")
	}
	if repo.order.TotalAmount != 350000 || response.TotalAmount != 350000 {
		t.Fatalf("order total = %d, response total = %d; want 350000", repo.order.TotalAmount, response.TotalAmount)
	}
	if len(repo.items) != 1 {
		t.Fatalf("saved order items = %d; want 1", len(repo.items))
	}
}
