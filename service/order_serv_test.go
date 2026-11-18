package service

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"ecommerce-catalog-api/domain"

	"github.com/google/uuid"
)

type orderRepositoryStub struct {
	order                    *domain.Order
	items                    []*domain.OrderItem
	paymentNotificationCalls int
}

func (r *orderRepositoryStub) CreateWithItems(_ context.Context, order *domain.Order, items []*domain.OrderItem) error {
	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = now
	order.StockReduced = false
	r.order = order
	r.items = append(r.items, items...)
	return nil
}

func (r *orderRepositoryStub) ProcessPaymentNotification(_ context.Context, externalID string, amount int, status domain.PaymentStatus, paymentMethod string) error {
	r.paymentNotificationCalls++
	if r.order.ExternalID == nil || *r.order.ExternalID != externalID || r.order.TotalAmount != amount {
		return fmt.Errorf("payment notification does not match order")
	}
	r.order.PaymentStatus = status
	if paymentMethod != "" {
		r.order.PaymentMethod = &paymentMethod
	}
	if status == domain.PaymentStatusSettled {
		r.order.StockReduced = true
	}
	return nil
}

func (r *orderRepositoryStub) FindByID(context.Context, uuid.UUID) (*domain.Order, error) {
	return nil, nil
}

func (r *orderRepositoryStub) FindByUserID(context.Context, uuid.UUID) ([]domain.Order, error) {
	return nil, nil
}

func (r *orderRepositoryStub) Update(_ context.Context, order *domain.Order) error {
	r.order = order
	return nil
}

func TestCreateOrderAllowsGuestCheckout(t *testing.T) {
	repo := &orderRepositoryStub{}
	snapServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/snap/v1/transactions" {
			t.Errorf("unexpected Midtrans request: %s %s", r.Method, r.URL.Path)
		}
		var request midtransTransactionRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("decode Midtrans request: %v", err)
		}
		if request.TransactionDetails.GrossAmount != 350000 {
			t.Errorf("Midtrans gross amount = %d; want 350000", request.TransactionDetails.GrossAmount)
		}
		if len(request.TransactionDetails.OrderID) > 50 || strings.Contains(request.TransactionDetails.OrderID, "/") {
			t.Errorf("Midtrans order ID is invalid: %q", request.TransactionDetails.OrderID)
		}
		if len(request.ItemDetails) != 1 || len(request.ItemDetails[0].Name) > 50 {
			t.Errorf("Midtrans item details have an invalid name: %+v", request.ItemDetails)
		}
		if r.Header.Get("Authorization") != "Basic a2V5LXRlc3Q6" {
			t.Errorf("unexpected Midtrans authorization: %q", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"token":"snap-token","redirect_url":"https://midtrans.test/snap"}`))
	}))
	defer snapServer.Close()

	service := NewOrderService(repo).(*orderService)
	service.midtrans = &midtransClient{
		serverKey:  "key-test",
		baseURL:    snapServer.URL,
		httpClient: snapServer.Client(),
	}
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

	response, err := service.CreateOrder(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateOrder() error = %v", err)
	}
	if repo.order == nil {
		t.Fatal("CreateOrder() did not save an order")
	}
	if repo.order.UserID != nil || response.UserID != nil {
		t.Fatal("guest order should not be associated with a user")
	}
	if repo.order.StockReduced || response.StockReduced {
		t.Fatal("stock should not be reduced before successful payment")
	}
	if repo.order.TotalAmount != 350000 || response.TotalAmount != 350000 {
		t.Fatalf("order total = %d, response total = %d; want 350000", repo.order.TotalAmount, response.TotalAmount)
	}
	if len(repo.items) != 1 {
		t.Fatalf("saved order items = %d; want 1", len(repo.items))
	}
	if response.SnapToken == nil || *response.SnapToken != "snap-token" {
		t.Fatal("checkout response should include the Midtrans Snap token")
	}
	if response.SnapRedirectUrl == nil || *response.SnapRedirectUrl != "https://midtrans.test/snap" {
		t.Fatal("checkout response should include the Midtrans redirect URL")
	}
}

func TestProcessPaymentWebhookReducesStockOnlyForSettledPayment(t *testing.T) {
	repo := &orderRepositoryStub{
		order: &domain.Order{
			ExternalID:    stringPointer("INV/20261117/0001"),
			TotalAmount:   350000,
			PaymentStatus: domain.PaymentStatusPending,
		},
	}
	service := NewOrderService(repo).(*orderService)
	service.midtrans = &midtransClient{serverKey: "key-test"}

	payload := validMidtransNotification("key-test", "INV/20261117/0001", "350000.00", "settlement")
	if err := service.ProcessPaymentWebhook(context.Background(), payload); err != nil {
		t.Fatalf("ProcessPaymentWebhook() error = %v", err)
	}
	if repo.order.PaymentStatus != domain.PaymentStatusSettled || !repo.order.StockReduced {
		t.Fatal("settled payment should update status and reduce stock")
	}

	repo.order.PaymentStatus = domain.PaymentStatusPending
	repo.order.StockReduced = false
	payload = validMidtransNotification("key-test", "INV/20261117/0001", "350000.00", "pending")
	if err := service.ProcessPaymentWebhook(context.Background(), payload); err != nil {
		t.Fatalf("ProcessPaymentWebhook(pending) error = %v", err)
	}
	if repo.order.PaymentStatus != domain.PaymentStatusPending || repo.order.StockReduced {
		t.Fatal("pending payment must not reduce stock")
	}
	if repo.paymentNotificationCalls != 1 {
		t.Fatalf("repository calls after settlement and pending notifications = %d; want 1", repo.paymentNotificationCalls)
	}
}

func TestProcessPaymentWebhookRejectsInvalidSignature(t *testing.T) {
	repo := &orderRepositoryStub{
		order: &domain.Order{
			ExternalID:  stringPointer("INV/20261117/0001"),
			TotalAmount: 350000,
		},
	}
	service := NewOrderService(repo).(*orderService)
	service.midtrans = &midtransClient{serverKey: "key-test"}

	payload := validMidtransNotification("key-test", "INV/20261117/0001", "350000.00", "settlement")
	payload["signature_key"] = "invalid"
	if err := service.ProcessPaymentWebhook(context.Background(), payload); err == nil {
		t.Fatal("ProcessPaymentWebhook() should reject an invalid signature")
	}
	if repo.order.StockReduced {
		t.Fatal("invalid notification must not reduce stock")
	}
}

func TestProcessPaymentWebhookAcknowledgesMidtransTestNotification(t *testing.T) {
	repo := &orderRepositoryStub{
		order: &domain.Order{
			ExternalID:    stringPointer("INV-20261118-08858ea9-0da0-4966-842b-4b353701a844"),
			TotalAmount:   105000,
			PaymentStatus: domain.PaymentStatusPending,
		},
	}
	service := NewOrderService(repo).(*orderService)
	service.midtrans = &midtransClient{serverKey: "key-test"}

	payload := validMidtransNotification("key-test", "payment_notif_test_M285664080_test-id", "105000.00", "settlement")
	payload["signature_key"] = "dashboard-test-signature"
	if err := service.ProcessPaymentWebhook(context.Background(), payload); err != nil {
		t.Fatalf("ProcessPaymentWebhook(test notification) error = %v", err)
	}
	if repo.order.PaymentStatus != domain.PaymentStatusPending || repo.order.StockReduced {
		t.Fatal("Midtrans test notification must not change a real order")
	}
}

func validMidtransNotification(serverKey, orderID, grossAmount, transactionStatus string) map[string]interface{} {
	statusCode := "200"
	hash := sha512.Sum512([]byte(orderID + statusCode + grossAmount + serverKey))
	return map[string]interface{}{
		"order_id":           orderID,
		"status_code":        statusCode,
		"gross_amount":       grossAmount,
		"signature_key":      hex.EncodeToString(hash[:]),
		"transaction_status": transactionStatus,
	}
}

func stringPointer(value string) *string {
	return &value
}
