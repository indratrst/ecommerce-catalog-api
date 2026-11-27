---
name: go-fiber-api
description: Architecture conventions, GORM transaction rules, and API standards for the Go Fiber e-commerce backend.
---

## Architecture & Layering
- Strictly follow `handler` -> `service` -> `repository` separation. Never put raw GORM queries inside handlers.
- Keep handler logic thin: parse params/body, validate input, call service, and return standardized JSON responses (`status`, `message`, `data`).

## Database & GORM Rules
- Database: PostgreSQL running via Docker.
- Wrap multi-table writes and inventory stock reductions in atomic `db.Transaction(func(tx *gorm.DB) error { ... })`.
- Always check `errors.Is(err, gorm.ErrRecordNotFound)` to return HTTP `404` instead of `500`.
- Support both authenticated (JWT) and guest checkout flows where applicable without breaking existing order schemas.

## Integrations & Webhooks
- Payments: Midtrans Snap. Always verify SHA-512 signature keys before mutating order/payment status.
- Shipping: RajaOngkir / Komerce tracking. Keep external HTTP client calls isolated in dedicated service helpers.

## Token-Saving Execution Rules
- Do not scan `vendor/` or large `.sum` files unless explicitly asked.
- After modifying `service` or `repository` files, run `go test ./...` or `go build ./...` to verify compilation before finishing.