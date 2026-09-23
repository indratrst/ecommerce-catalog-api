package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleAdmin    UserRole = "admin"
)

type User struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	FullName     string    `json:"full_name" gorm:"type:varchar(100);not null"`
	Email        string    `json:"email" gorm:"type:varchar(150);not null;unique"`
	PasswordHash string    `json:"-" gorm:"type:varchar(255);not null"` // "-" agar password hash tidak di-serialize ke JSON
	PhoneNumber  *string   `json:"phone_number,omitempty" gorm:"type:varchar(20)"`
	Role         UserRole  `json:"role" gorm:"type:user_role;default:'customer'"`
	Addresses    []Address `json:"addresses,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Address struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID        uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	RecipientName string    `json:"recipient_name" gorm:"type:varchar(100);not null"`
	PhoneNumber   string    `json:"phone_number" gorm:"type:varchar(20);not null"`
	StreetAddress string    `json:"street_address" gorm:"type:text;not null"`
	CityID        int       `json:"city_id" gorm:"not null"`
	CityName      string    `json:"city_name" gorm:"type:varchar(100);not null"`
	ProvinceID    int       `json:"province_id" gorm:"not null"`
	ProvinceName  string    `json:"province_name" gorm:"type:varchar(100);not null"`
	PostalCode    string    `json:"postal_code" gorm:"type:varchar(10);not null"`
	IsPrimary     bool      `json:"is_primary" gorm:"default:false"`
	CreatedAt     time.Time `json:"created_at"`
}
