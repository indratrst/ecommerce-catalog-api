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

	// 2. Enum role
	db.Exec(`
    DO $$ BEGIN 
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role_enum') THEN 
            CREATE TYPE role_enum AS ENUM ('SUPERUSER', 'ADMIN', 'USER');
        END IF; 
    END $$;
`)

	// 3. Enum shipping_method
	db.Exec(`
    DO $$ BEGIN 
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'shipping_method_enum') THEN 
            CREATE TYPE shipping_method_enum AS ENUM ('PICKUP_STORE', 'SHIPPING');
        END IF; 
    END $$;
`)

	// 4. Enum payment_status
	db.Exec(`
    DO $$ BEGIN 
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'payment_status_enum') THEN 
            CREATE TYPE payment_status_enum AS ENUM ('PENDING', 'SETTLEMENT', 'EXPIRED', 'CANCEL', 'FAILED');
        END IF; 
    END $$;
`)

	// 5. Enum fulfillment_status
	db.Exec(`
    DO $$ BEGIN 
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'fulfillment_status_enum') THEN 
            CREATE TYPE fulfillment_status_enum AS ENUM ('NOT_APPLICABLE', 'PROCESSING', 'READY_TO_PICKUP', 'PICKED_UP', 'SHIPPED', 'DELIVERED');
        END IF; 
    END $$;
`)
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Category{},
		&domain.Product{},
		&domain.ProductVariant{},
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
