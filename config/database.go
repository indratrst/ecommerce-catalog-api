package config

import (
	"ecommerce-catalog-api/domain"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "rahasia")
	dbName := getEnv("DB_NAME", "ecommerce_catalog_db")
	port := getEnv("DB_PORT", "5432")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbName, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database: ", err)
	}
	// Setup Connection Pooling
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Gagal mendapatkan instance sql.DB: ", err)
	}

	sqlDB.SetMaxIdleConns(10)               // Jumlah koneksi idle maksimum
	sqlDB.SetMaxOpenConns(100)              // Jumlah koneksi terbuka maksimum
	sqlDB.SetConnMaxLifetime(1 * time.Hour) // Umur maksimum koneksi

	fmt.Println("Berhasil terhubung ke PostgreSQL!")

	// 1. Extension UUID
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`)

	// 2. Enum user_role
	db.Exec(`
    DO $$ BEGIN 
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN 
            CREATE TYPE user_role AS ENUM ('customer', 'admin'); 
        END IF; 
    END $$;
`)

	// 3. Enum coupon_type
	db.Exec(`
    DO $$ BEGIN 
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'coupon_type') THEN 
            CREATE TYPE coupon_type AS ENUM ('percentage', 'fixed'); 
        END IF; 
    END $$;
`)

	// 4. Enum delivery_type (Omnichannel support: courier vs pickup in store)
	db.Exec(`
    DO $$ BEGIN 
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'delivery_type') THEN 
            CREATE TYPE delivery_type AS ENUM ('courier', 'pickup'); 
        END IF; 
    END $$;
`)

	// 5. Enum order_status
	db.Exec(`
    DO $$ BEGIN 
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status') THEN 
            CREATE TYPE order_status AS ENUM ('pending_payment', 'paid', 'processing', 'shipped', 'ready_for_pickup', 'completed', 'cancelled'); 
        END IF; 
    END $$;
`)
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Address{},
		&domain.Category{},
		&domain.Size{},
		&domain.Color{},
		&domain.Product{},
		&domain.ProductVariant{},
		&domain.ProductImage{},
		&domain.Wishlist{},
		&domain.Review{},
		&domain.CartItem{},
		&domain.Coupon{},
		&domain.Order{},
		&domain.OrderItem{},
	)
	if err != nil {
		log.Fatal("Gagal melakukan AutoMigrate: ", err)
	}

	return db
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
