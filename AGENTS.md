# AGENTS.md

- Go 1.26 module `ecommerce-catalog-api`: Fiber v2 + GORM/PostgreSQL REST API. `main.go` is the single entrypoint and wires every layer manually; there is no `cmd/`, `internal/`, monorepo, Makefile, CI, lint, or codegen config despite what `Ecommerce_docs/Structure_FolderBE.MD` describes.
- `main.go` is the route source of truth. `README.md` and `Ecommerce_docs/Api_contract.MD` are stale: the real auth routes are `/api/v1/auth/register` and `/api/v1/auth/login`, and the documented cart/coupon endpoints are not implemented (only domain models exist).
- Keep `handler/` → `service/` → `repository/` separation with contracts in `domain/`; return all HTTP responses through `pkg/response.Success/Error` (`success`, `message`, `data`/`errors`). Also follow `.opencode/skills/go-fiber-api/SKILL.md` for GORM transaction and webhook rules.

## Verify

- `go build ./...` passes. `go test ./...` has a pre-existing failure: `TestProcessPaymentWebhookRejectsInvalidSignature` builds a `midtransClient` without `httpClient`, so an invalid signature enters the Midtrans status-API fallback and panics. Don't attribute it to unrelated edits. `go test ./handler` and the other `service` tests pass.
- `docker compose up --build -d` is required after Go changes (the image does not hot-reload); logs: `docker compose logs -f api`. Host PostgreSQL port is 5433; inside Compose the DB host is `postgres`.

## Runtime gotchas

- `config.ConnectDB()` creates the `pgcrypto` extension, enum types, runs GORM `AutoMigrate`, and adds `image` columns on every startup. There are no migration files; restart the API against the intended DB to apply schema changes.
- Compose loads `.env` via `env_file`, but its `environment:` block overrides DB/JWT/UPLOAD values and the `MIDTRANS_*` passthrough lines are commented out. Host-exported Midtrans variables do not reach the API; edit `.env` instead. No `.env.example` is tracked.
- `.ignore` hides `.env` and `go.sum` from ripgrep; read them explicitly when needed.
- Uploads: `POST /api/v1/uploads` (ADMIN/SUPERUSER) takes multipart `file`, max 5 MB JPEG/PNG/GIF/WEBP. `UPLOAD_DIR` defaults to `uploads/` and must be persistent outside Compose (Compose mounts the `api_uploads` volume); `PUBLIC_API_URL` sets the host in returned URLs.
- Roles are uppercase `SUPERUSER`, `ADMIN`, `USER`. User management requires `SUPERUSER`; uploads and shipment updates require ADMIN or SUPERUSER.
- Checkout uses `OptionalProtected()` so guests are allowed. Midtrans item total including shipping must equal `order.TotalAmount`; stock decrements only on settled/success/captured-with-accept webhooks, and `payment_notif_test_*` IDs are acknowledged without changing orders.
