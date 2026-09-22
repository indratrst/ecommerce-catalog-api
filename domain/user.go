package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name        string `json:"name"`
	Email       string `json:"email" gorm:"unique;not null"`
	MobilePhone string `json:"mobile_phone"`
	Password    string `json:"-"` // Hidden dari JSON response
	Token       string `json:"token,omitempty" gorm:"type:text"`
}

type RegisterRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	MobilePhone string `json:"mobile_phone"`
	Password    string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// Contract Repository
type UserRepository interface {
	Create(user *User) error
	Update(user *User) error
	FindByEmail(email string) (User, error)
	FindByID(id uint) (User, error)
}

// Contract Service
type AuthService interface {
	Register(req RegisterRequest) (AuthResponse, error)
	Login(req LoginRequest) (AuthResponse, error)
}
