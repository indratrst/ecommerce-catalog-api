Berikut adalah rancangan ERD (Entity-Relationship Diagram) PostgreSQL lengkap yang disesuaikan secara presisi dengan seluruh spesifikasi pada PRD V1.

Diagram Relasi Ringkasan (High-Level Overview)
```text
               ┌──────────────┐
               │   Category   │
               └──────┬───────┘
                      │ 1
                      │
                      │ *
┌──────────────┐ 1   *┌──────────────┐ 1     *┌─────────────────┐
│     User     ├──────┤    Product   ├───────┤ ProductVariant  │
└──────┬───────┘      └──────────────┘        └────────┬────────┘
       │ 1                                             │ 1
       │                                               │
       │ *                                             │ *
┌──────┴───────┐                             ┌─────────┴───────┐
│    Order     ├────────────────────────────►│    OrderItem    │
└──────────────┘ 1                         * └─────────────────┘

```



# DDL SQL Complete Schema (PostgreSQL)
-- Extension untuk UUID generator

```text
-- Extension untuk UUID generator
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =============================================================================
-- ENUM TYPES
-- =============================================================================

CREATE TYPE role_enum AS ENUM ('SUPERUSER', 'ADMIN', 'USER');

CREATE TYPE shipping_method_enum AS ENUM ('PICKUP_STORE', 'SHIPPING');

CREATE TYPE payment_status_enum AS ENUM ('PENDING', 'SETTLEMENT', 'EXPIRED', 'CANCEL', 'FAILED');

CREATE TYPE fulfillment_status_enum AS ENUM (
    'NOT_APPLICABLE',
    'PROCESSING',
    'READY_TO_PICKUP',
    'PICKED_UP',
    'SHIPPED',
    'DELIVERED'
);

-- =============================================================================
-- 1. USER MODULE
-- =============================================================================

CREATE TABLE "User" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) NOT NULL UNIQUE,
    "password" VARCHAR(255) NOT NULL,
    "role" role_enum NOT NULL DEFAULT 'USER',
    "image" TEXT,
    "createdAt" TIMESTAMP(3) WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- 2. PRODUCT & CATEGORY MODULE
-- =============================================================================

CREATE TABLE "Category" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" VARCHAR(255) NOT NULL,
    "slug" VARCHAR(255) NOT NULL UNIQUE,
    "description" TEXT,
    "image" TEXT,
    "createdAt" TIMESTAMP(3) WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "Product" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "categoryId" UUID NOT NULL REFERENCES "Category"("id") ON DELETE RESTRICT ON UPDATE CASCADE,
    "title" VARCHAR(255) NOT NULL,
    "price" INT NOT NULL,
    "originalPrice" INT,
    "description" TEXT NOT NULL,
    "image" TEXT,
    "ratingRate" DOUBLE PRECISION DEFAULT 0,
    "ratingCount" INT DEFAULT 0,
    "createdAt" TIMESTAMP(3) WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "ProductVariant" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "productId" UUID NOT NULL REFERENCES "Product"("id") ON DELETE CASCADE ON UPDATE CASCADE,
    "size" VARCHAR(50) NOT NULL,
    "color" VARCHAR(50),
    "stock" INT NOT NULL DEFAULT 0,
    "isActive" BOOLEAN NOT NULL DEFAULT TRUE,
    "createdAt" TIMESTAMP(3) WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- 3. ORDER MODULE
-- =============================================================================

CREATE TABLE "Order" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "userId" UUID REFERENCES "User"("id") ON DELETE SET NULL ON UPDATE CASCADE,
    "customerName" VARCHAR(255) NOT NULL,
    "customerEmail" VARCHAR(255) NOT NULL,
    "customerPhone" VARCHAR(50) NOT NULL,
    "shippingAddress" TEXT,
    "shippingMethod" shipping_method_enum NOT NULL,
    "shippingCourier" VARCHAR(100),
    "trackingNumber" VARCHAR(100),
    "totalAmount" INT NOT NULL,
    "paymentStatus" payment_status_enum NOT NULL DEFAULT 'PENDING',
    "paymentMethod" VARCHAR(100),
    "snapToken" VARCHAR(255),
    "snapRedirectUrl" TEXT,
    "externalId" VARCHAR(255) UNIQUE,
    "stockReduced" BOOLEAN NOT NULL DEFAULT FALSE,
    "fulfillmentStatus" fulfillment_status_enum NOT NULL DEFAULT 'NOT_APPLICABLE',
    "createdAt" TIMESTAMP(3) WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "OrderItem" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "orderId" UUID NOT NULL REFERENCES "Order"("id") ON DELETE CASCADE ON UPDATE CASCADE,
    "productVariantId" UUID NOT NULL REFERENCES "ProductVariant"("id") ON DELETE RESTRICT ON UPDATE CASCADE,
    "quantity" INT NOT NULL,
    "priceAtPurchase" INT NOT NULL,
    "variantSnapshot" JSONB
);

-- =============================================================================
-- 4. INDEXES
-- =============================================================================

CREATE INDEX "idx_product_id" ON "Product"("id");
CREATE INDEX "idx_product_variant_product_id" ON "ProductVariant"("productId");
CREATE INDEX "idx_order_user_id" ON "Order"("userId");
CREATE INDEX "idx_order_external_id" ON "Order"("externalId");
CREATE INDEX "idx_order_item_order_id" ON "OrderItem"("orderId");
CREATE INDEX "idx_order_item_variant_id" ON "OrderItem"("productVariantId");
```