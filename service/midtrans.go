package service

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"ecommerce-catalog-api/domain"
)

const (
	midtransSandboxURL    = "https://app.sandbox.midtrans.com"
	midtransProductionURL = "https://app.midtrans.com"
)

type midtransClient struct {
	serverKey  string
	baseURL    string
	httpClient *http.Client
}

type midtransTransactionRequest struct {
	TransactionDetails struct {
		OrderID     string `json:"order_id"`
		GrossAmount int    `json:"gross_amount"`
	} `json:"transaction_details"`
	CustomerDetails struct {
		FirstName string `json:"first_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	} `json:"customer_details"`
	ItemDetails []midtransItem `json:"item_details"`
}

type midtransItem struct {
	ID       string `json:"id"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
	Name     string `json:"name"`
}

type midtransTransactionResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

func newMidtransClientFromEnv() *midtransClient {
	baseURL := midtransSandboxURL
	if strings.EqualFold(os.Getenv("MIDTRANS_IS_PRODUCTION"), "true") {
		baseURL = midtransProductionURL
	}

	return &midtransClient{
		serverKey: os.Getenv("MIDTRANS_SERVER_KEY"),
		baseURL:   baseURL,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *midtransClient) createTransaction(ctx context.Context, order *domain.Order, items []*domain.OrderItem) (*midtransTransactionResponse, error) {
	if strings.TrimSpace(c.serverKey) == "" {
		return nil, fmt.Errorf("MIDTRANS_SERVER_KEY is required")
	}

	requestBody := midtransTransactionRequest{}
	if order.ExternalID == nil || strings.TrimSpace(*order.ExternalID) == "" {
		return nil, fmt.Errorf("order externalId is required for Midtrans")
	}
	requestBody.TransactionDetails.OrderID = *order.ExternalID
	requestBody.CustomerDetails.FirstName = order.CustomerName
	requestBody.CustomerDetails.Email = order.CustomerEmail
	requestBody.CustomerDetails.Phone = order.CustomerPhone
	var itemDetailsTotal int64
	for _, item := range items {
		if item.PriceAtPurchase < 0 || item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid Midtrans item price or quantity")
		}
		quantity := int64(item.Quantity)
		price := int64(item.PriceAtPurchase)
		if price > 0 && quantity > int64(^uint64(0)>>1)/price {
			return nil, fmt.Errorf("Midtrans item subtotal exceeds supported amount")
		}
		subtotal := price * quantity
		if itemDetailsTotal > int64(^uint64(0)>>1)-subtotal {
			return nil, fmt.Errorf("Midtrans item total exceeds supported amount")
		}
		itemDetailsTotal += subtotal

		requestBody.ItemDetails = append(requestBody.ItemDetails, midtransItem{
			ID:       item.ProductVariantID.String(),
			Price:    item.PriceAtPurchase,
			Quantity: item.Quantity,
			Name:     "Variant " + item.ProductVariantID.String(),
		})
	}
	if itemDetailsTotal != int64(order.TotalAmount) {
		return nil, fmt.Errorf("order total %d does not match Midtrans item total %d", order.TotalAmount, itemDetailsTotal)
	}
	requestBody.TransactionDetails.GrossAmount = int(itemDetailsTotal)
	if int64(requestBody.TransactionDetails.GrossAmount) != itemDetailsTotal {
		return nil, fmt.Errorf("Midtrans total exceeds supported amount")
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to encode Midtrans transaction: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(c.baseURL, "/")+"/snap/v1/transactions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create Midtrans request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.serverKey+":")))

	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create Midtrans transaction: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("failed to read Midtrans response: %w", err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("Midtrans transaction request failed with status %d: %s", response.StatusCode, strings.TrimSpace(string(responseBody)))
	}

	var transaction midtransTransactionResponse
	if err := json.Unmarshal(responseBody, &transaction); err != nil {
		return nil, fmt.Errorf("failed to decode Midtrans response: %w", err)
	}
	if transaction.Token == "" || transaction.RedirectURL == "" {
		return nil, fmt.Errorf("Midtrans response is missing token or redirect_url")
	}
	return &transaction, nil
}

func (c *midtransClient) verifyNotification(payload map[string]interface{}) (string, int, string, domain.PaymentStatus, error) {
	if strings.TrimSpace(c.serverKey) == "" {
		return "", 0, "", "", fmt.Errorf("MIDTRANS_SERVER_KEY is required")
	}
	orderID, err := requiredPayloadString(payload, "order_id")
	if err != nil {
		return "", 0, "", "", err
	}
	statusCode, err := requiredPayloadString(payload, "status_code")
	if err != nil {
		return "", 0, "", "", err
	}
	grossAmount, err := requiredPayloadString(payload, "gross_amount")
	if err != nil {
		return "", 0, "", "", err
	}
	signature, err := requiredPayloadString(payload, "signature_key")
	if err != nil {
		return "", 0, "", "", err
	}
	transactionStatus, err := requiredPayloadString(payload, "transaction_status")
	if err != nil {
		return "", 0, "", "", err
	}

	hash := sha512.Sum512([]byte(orderID + statusCode + grossAmount + c.serverKey))
	expectedSignature := hex.EncodeToString(hash[:])
	if !equalSignature(signature, expectedSignature) {
		return "", 0, "", "", fmt.Errorf("invalid Midtrans signature")
	}

	amount, err := parseMidtransAmount(grossAmount)
	if err != nil {
		return "", 0, "", "", err
	}
	paymentMethod, _ := payload["payment_type"].(string)
	fraudStatus, _ := payload["fraud_status"].(string)
	status, err := mapMidtransPaymentStatus(transactionStatus, fraudStatus)
	if err != nil {
		return "", 0, "", "", err
	}
	return orderID, amount, paymentMethod, status, nil
}

func requiredPayloadString(payload map[string]interface{}, key string) (string, error) {
	value, ok := payload[key].(string)
	if !ok || strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("Midtrans notification field %q is required", key)
	}
	return value, nil
}

func equalSignature(actual, expected string) bool {
	actualBytes, actualErr := hex.DecodeString(actual)
	expectedBytes, expectedErr := hex.DecodeString(expected)
	return actualErr == nil && expectedErr == nil && hmac.Equal(actualBytes, expectedBytes)
}

func parseMidtransAmount(value string) (int, error) {
	amount, ok := new(big.Rat).SetString(value)
	if !ok || !amount.IsInt() || !amount.Num().IsInt64() {
		return 0, fmt.Errorf("invalid Midtrans gross_amount %q", value)
	}
	amountInt64 := amount.Num().Int64()
	parsedAmount := int(amountInt64)
	if int64(parsedAmount) != amountInt64 || parsedAmount < 0 {
		return 0, fmt.Errorf("invalid Midtrans gross_amount %q", value)
	}
	return parsedAmount, nil
}

func mapMidtransPaymentStatus(transactionStatus, fraudStatus string) (domain.PaymentStatus, error) {
	switch strings.ToLower(strings.TrimSpace(transactionStatus)) {
	case "settlement", "success":
		return domain.PaymentStatusSettled, nil
	case "capture":
		if strings.EqualFold(strings.TrimSpace(fraudStatus), "accept") {
			return domain.PaymentStatusSettled, nil
		}
		return domain.PaymentStatusPending, nil
	case "pending":
		return domain.PaymentStatusPending, nil
	case "expire":
		return domain.PaymentStatusExpired, nil
	case "cancel":
		return domain.PaymentStatusCancel, nil
	case "deny", "failure":
		return domain.PaymentStatusFailed, nil
	default:
		return "", fmt.Errorf("unsupported Midtrans transaction_status %q", transactionStatus)
	}
}
