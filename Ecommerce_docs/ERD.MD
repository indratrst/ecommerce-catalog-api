Berikut adalah rancangan ERD (Entity-Relationship Diagram) PostgreSQL lengkap yang disesuaikan secara presisi dengan seluruh spesifikasi pada PRD V1.

Diagram Relasi Ringkasan (High-Level Overview)
```text
[users] 1───* [addresses]
  │   1───* [orders] 1───* [order_items] *───1 [product_variants]
  │   1───* [cart_items] *───1 [product_variants]
  │   1───* [wishlists]  *───1 [products]
  │   1───* [reviews]    *───1 [products]
  │
[categories] 1───* [products] 1───* [product_variants] *───1 [colors]
                         │                   │
                         │                   *───1 [sizes]
                         │
                         *───* [product_images]

[coupons] 1───* [orders]

```



# DDL SQL Complete Schema (PostgreSQL)
-- Extension untuk UUID generator

```text
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =============================================================================
-- 1. USER & AUTHENTICATION MODULE
-- =============================================================================

CREATE TYPE user_role AS ENUM ('customer', 'admin');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name VARCHAR(100) NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20),
    role user_role DEFAULT 'customer',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE addresses (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recipient_name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    street_address TEXT NOT NULL,
    city_id INT NOT NULL, -- ID Kota dari API Ekspedisi (RajaOngkir/Biteship)
    city_name VARCHAR(100) NOT NULL,
    province_id INT NOT NULL,
    province_name VARCHAR(100) NOT NULL,
    postal_code VARCHAR(10) NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- 2. PRODUCT & INVENTORY MODULE
-- =============================================================================

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    parent_id INT REFERENCES categories(id) ON DELETE SET NULL, -- Memungkinkan Sub-Kategori
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sizes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) NOT NULL, -- Contoh: S, M, L, XL, EU 42
    category_type VARCHAR(50) -- Contoh: Apparel, Footwear
);

CREATE TABLE colors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,  -- Contoh: Navy Blue
    hex_code VARCHAR(7) NOT NULL  -- Contoh: #000080 (UI Swatch di Next.js)
);

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id INT REFERENCES categories(id) ON DELETE RESTRICT,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    base_price DECIMAL(12, 2) NOT NULL,
    rating_avg DECIMAL(3, 2) DEFAULT 0.0,
    review_count INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE product_variants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    size_id INT NOT NULL REFERENCES sizes(id) ON DELETE RESTRICT,
    color_id INT REFERENCES colors(id) ON DELETE RESTRICT, -- Opsional sesuai PRD
    sku VARCHAR(100) NOT NULL UNIQUE,
    price_override DECIMAL(12, 2), -- Jika harga bervariasi per ukuran/warna
    stock_quantity INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_variant UNIQUE (product_id, size_id, color_id)
);

CREATE TABLE product_images (
    id SERIAL PRIMARY KEY,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE,
    display_order INT DEFAULT 0
);

-- =============================================================================
-- 3. WISHLIST & REVIEWS MODULE
-- =============================================================================

CREATE TABLE wishlists (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_user_wishlist UNIQUE (user_id, product_id)
);

CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    rating INT CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_user_review UNIQUE (user_id, product_id)
);

-- =============================================================================
-- 4. CART & MARKETING MODULE
-- =============================================================================

CREATE TABLE cart_items (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    variant_id UUID NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
    quantity INT NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_user_cart_variant UNIQUE (user_id, variant_id)
);

CREATE TYPE coupon_type AS ENUM ('percentage', 'fixed_amount');

CREATE TABLE coupons (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    discount_type coupon_type NOT NULL,
    discount_value DECIMAL(12, 2) NOT NULL,
    min_purchase DECIMAL(12, 2) DEFAULT 0.0,
    max_discount DECIMAL(12, 2), -- Maksimal potongan jika %
    expired_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- 5. ORDER & PAYMENT MODULE
-- =============================================================================

CREATE TYPE order_status AS ENUM ('pending_payment', 'paid', 'processing', 'shipped', 'completed', 'cancelled');

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number VARCHAR(50) NOT NULL UNIQUE, -- Contoh: INV/2026/09/001
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    coupon_id INT REFERENCES coupons(id) ON DELETE SET NULL,
    
    -- Snapshots Informasi Alamat Pengiriman
    shipping_recipient VARCHAR(100) NOT NULL,
    shipping_phone VARCHAR(20) NOT NULL,
    shipping_address TEXT NOT NULL,
    shipping_courier VARCHAR(50) NOT NULL, -- JNE/TIKI/Pos
    shipping_service VARCHAR(50) NOT NULL, -- REG/YES
    tracking_number VARCHAR(100),          -- Nomor Resi

    -- Rincian Biaya
    subtotal DECIMAL(12, 2) NOT NULL,
    discount_amount DECIMAL(12, 2) DEFAULT 0.0,
    shipping_cost DECIMAL(12, 2) NOT NULL,
    total_amount DECIMAL(12, 2) NOT NULL,

    -- Midtrans / Xendit
    payment_token VARCHAR(255),
    payment_url TEXT,
    status order_status DEFAULT 'pending_payment',

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    variant_id UUID NOT NULL REFERENCES product_variants(id) ON DELETE RESTRICT,
    
    -- Snapshot data produk saat checkout dilakukan
    product_name VARCHAR(255) NOT NULL,
    variant_size VARCHAR(20) NOT NULL,
    variant_color VARCHAR(50),
    price DECIMAL(12, 2) NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    subtotal DECIMAL(12, 2) NOT NULL
);

-- =============================================================================
-- 6. INDEXES FOR PERFORMANCE OPTIMIZATION
-- =============================================================================

CREATE INDEX idx_products_slug ON products(slug);
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_variants_product ON product_variants(product_id);
CREATE INDEX idx_orders_user ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
```