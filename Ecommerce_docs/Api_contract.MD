# REST API Specification — Fashion E-Commerce V1

Dokumen ini mendefinisikan *API Contract* antara **Backend Go (Fiber)** dan **Frontend Next.js (App Router)** untuk platform e-commerce fashion.

- **Base URL:** `/api/v1`
- **Content-Type:** `application/json`
- **Authentication:** `Bearer <JWT_TOKEN>` pada HTTP Header `Authorization`

---

## Standard Response Formats

### 1. Success Response
```json
{
  "success": true,
  "message": "Data retrieved successfully",
  "data": {}
}
```

### 2. Error Response
```json
{
  "success": false,
  "message": "Validation error occurred",
  "errors": [
    {
      "field": "email",
      "message": "must be a valid email address"
    }
  ]
}
```

### 3. Paginated Response Structure
```json
{
  "success": true,
  "message": "Products retrieved successfully",
  "data": {
    "items": [],
    "pagination": {
      "current_page": 1,
      "total_pages": 5,
      "total_items": 50,
      "limit": 10
    }
  }
}
```

## Detailed Endpoint Specifications
### 1. Authentication & Users
POST /auth/register
Public — Registrasi akun customer baru.

Request Body:

```JSON
{
  "full_name": "Indra Tristia",
  "email": "indra@example.com",
  "password": "Password123!",
  "phone_number": "081234567890"
}
```
Response Success (201 Created):

```JSON
{
  "success": true,
  "message": "Registration successful",
  "data": {
    "user_id": "e9b1a8d0-2f3c-4e5a-8b1c-3d2e1f0a9b8c",
    "email": "indra@example.com"
  }
}
```

POST /auth/login
Public — Login dan mendapatkan JWT Token.

Request Body:

```JSON
{
  "email": "indra@example.com",
  "password": "Password123!"
}


Response Success (200 OK):

```JSON
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "e9b1a8d0-2f3c-4e5a-8b1c-3d2e1f0a9b8c",
      "full_name": "Indra Tristia",
      "email": "indra@example.com",
      "role": "customer"
    }
  }
}
```

## 2. Catalog & Discovery
GET /products
Public — Mengambil daftar katalog produk dengan opsi pencarian, pemfilteran, dan pengurutan.

Query Parameters:

```text

search (string, optional): Kata kunci judul produk.

category_id (int, optional): ID kategori produk.

price_min (number, optional): Filter harga minimum.

price_max (number, optional): Filter harga maksimum.

sort (string, optional): price_asc, price_desc, newest (default: newest).

page (int, default: 1): Nomor halaman.

limit (int, default: 10): Jumlah item per halaman.

```

Response Success (200 OK):

```JSON
{
  "success": true,
  "message": "Products retrieved successfully",
  "data": {
    "items": [
      {
        "id": "e9b1a8d0-2f3c-4e5a-8b1c-3d2e1f0a9b8c",
        "title": "Kemeja Casual Slimfit Navy",
        "slug": "kemeja-casual-slimfit-navy",
        "base_price": 249000,
        "rating_avg": 4.8,
        "review_count": 24,
        "primary_image": "[https://storage.example.com/products/kemeja-navy-1.jpg](https://storage.example.com/products/kemeja-navy-1.jpg)",
        "category": {
          "id": 2,
          "name": "Kemeja Pria"
        }
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 5,
      "total_items": 48,
      "limit": 10
    }
  }
}
```

GET /products/:slug
Public — Mengambil detail lengkap produk untuk halaman PDP (Product Detail Page).

Response Success (200 OK):

``` JSON
{
  "success": true,
  "message": "Product detail retrieved successfully",
  "data": {
    "id": "e9b1a8d0-2f3c-4e5a-8b1c-3d2e1f0a9b8c",
    "title": "Kemeja Casual Slimfit Navy",
    "slug": "kemeja-casual-slimfit-navy",
    "description": "Kemeja berbahan katun poplin premium yang adem dan nyaman digunakan sehari-hari.",
    "base_price": 249000,
    "rating_avg": 4.8,
    "review_count": 24,
    "category": {
      "id": 2,
      "name": "Kemeja Pria"
    },
    "images": [
      {
        "id": 101,
        "image_url": "[https://storage.example.com/products/kemeja-navy-1.jpg](https://storage.example.com/products/kemeja-navy-1.jpg)",
        "is_primary": true
      },
      {
        "id": 102,
        "image_url": "[https://storage.example.com/products/kemeja-navy-2.jpg](https://storage.example.com/products/kemeja-navy-2.jpg)",
        "is_primary": false
      }
    ],
    "variants": [
      {
        "id": "c1f2e3d4-5a6b-7c8d-9e0f-1a2b3c4d5e6f",
        "sku": "KMJ-NVY-M",
        "price_override": null,
        "effective_price": 249000,
        "stock_quantity": 15,
        "size": { "id": 1, "name": "M" },
        "color": { "id": 5, "name": "Navy Blue", "hex_code": "#000080" }
      },
      {
        "id": "d2f3e4d5-6a7b-8c9d-0e1f-2a3b4c5d6e7f",
        "sku": "KMJ-NVY-L",
        "price_override": 259000,
        "effective_price": 259000,
        "stock_quantity": 8,
        "size": { "id": 2, "name": "L" },
        "color": { "id": 5, "name": "Navy Blue", "hex_code": "#000080" }
      }
    ]
  }
}
```

## 3. Cart & Coupon
POST /cart
JWT Required — Menambahkan varian produk ke keranjang belanja.

Request Body:

```JSON
{
  "variant_id": "c1f2e3d4-5a6b-7c8d-9e0f-1a2b3c4d5e6f",
  "quantity": 2
}

```

Response Success (200 OK / 201 Created):

```JSON
{
  "success": true,
  "message": "Item added to cart successfully",
  "data": {
    "cart_item_id": 15,
    "variant_id": "c1f2e3d4-5a6b-7c8d-9e0f-1a2b3c4d5e6f",
    "quantity": 2
  }
}
```

POST /coupons/validate
JWT Required — Memverifikasi dan menghitung besaran diskon kode kupon.

Request Body:

```JSON
{
  "code": "DISMARET10",
  "subtotal": 498000
}
```

Response Success (200 OK):

```JSON
{
  "success": true,
  "message": "Coupon applied successfully",
  "data": {
    "coupon_id": 3,
    "code": "DISMARET10",
    "discount_type": "percentage",
    "discount_amount": 49800
  }
}
```

## 4. Orders & Checkout
POST /orders/checkout
JWT Required — Membuat pesanan baru, memesan stok, dan mendapatkan token Payment Gateway (Midtrans/Xendit).

Request Body:

```JSON
{
  "address_id": 12,
  "courier": "jne",
  "courier_service": "REG",
  "coupon_code": "DISMARET10",
  "items": [
    {
      "variant_id": "c1f2e3d4-5a6b-7c8d-9e0f-1a2b3c4d5e6f",
      "quantity": 2
    }
  ]
}
```
Response Success (201 Created):

```JSON
{
  "success": true,
  "message": "Order created successfully",
  "data": {
    "order_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "order_number": "INV/20260922/0001",
    "subtotal": 498000,
    "discount_amount": 49800,
    "shipping_cost": 22000,
    "total_amount": 470200,
    "payment_token": "SNAP-TOKEN-MIDTRANS-123456",
    "payment_url": "[https://app.sandbox.midtrans.com/snap/v2/vtweb/SNAP-TOKEN-MIDTRANS-123456](https://app.sandbox.midtrans.com/snap/v2/vtweb/SNAP-TOKEN-MIDTRANS-123456)"
  }
}
```
## 5. Webhooks
POST /webhooks/payment
Signature Verified — Menerima notifikasi server-to-server dari Midtrans/Xendit untuk memperbarui status pesanan menjadi paid atau cancelled.

Request Body (Midtrans Payload Example):

```JSON
{
  "transaction_time": "2026-09-22 15:30:00",
  "transaction_status": "settlement",
  "transaction_id": "827364-abcd-1234",
  "status_message": "midtrans payment status is settlement",
  "order_id": "INV/20260922/0001",
  "gross_amount": "470200.00",
  "signature_key": "a1b2c3d4e5f6..."
}
```

Response Success (200 OK):

```JSON
{
  "success": true,
  "message": "Webhook processed successfully"
}
```

## 6. Admin Operations
PATCH /admin/orders/:id/shipment
Admin JWT Required — Memperbarui nomor resi pengiriman dan status pesanan menjadi shipped.

Request Body:

```JSON
{
  "tracking_number": "JX123456789ID",
  "status": "shipped"
}
```
Response Success (200 OK):

```JSON
{
  "success": true,
  "message": "Shipment updated successfully",
  "data": {
    "order_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "tracking_number": "JX123456789ID",
    "status": "shipped"
  }
}