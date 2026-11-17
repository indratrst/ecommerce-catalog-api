package handler

import (
	"testing"
)

func TestParsePaymentWebhookPayload(t *testing.T) {
	payload, err := parsePaymentWebhookPayload([]byte(`{"order_id":"INV-1","transaction_status":"settlement"}`))
	if err != nil {
		t.Fatalf("parsePaymentWebhookPayload() error = %v", err)
	}
	if payload["order_id"] != "INV-1" {
		t.Fatalf("order_id = %v; want INV-1", payload["order_id"])
	}
}

func TestParsePaymentWebhookPayloadRejectsInvalidJSON(t *testing.T) {
	if _, err := parsePaymentWebhookPayload([]byte(`{"order_id":`)); err == nil {
		t.Fatal("parsePaymentWebhookPayload() should reject invalid JSON")
	}
}

func TestParsePaymentWebhookPayloadRejectsEmptyBody(t *testing.T) {
	if _, err := parsePaymentWebhookPayload(nil); err == nil || err.Error() != "request body is empty; expected a Midtrans JSON notification" {
		t.Fatalf("parsePaymentWebhookPayload() error = %v; want explicit empty body error", err)
	}
}
