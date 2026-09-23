package domain

import "context"

// 1. Interface AuthService (Contract untuk Business Logic Auth)
type AuthService interface {
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error) // sesuaikan tipe data ID (string/uint)
	Update(ctx context.Context, user *User) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type RegisterRequest struct {
	FullName    string `json:"full_name" validate:"required,min=3,max=100"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	PhoneNumber string `json:"phone_number" validate:"required,numeric,min=10,max=15"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID          string   `json:"id"`
	FullName    string   `json:"full_name"`
	PhoneNumber *string  `json:"phone_number,omitempty"`
	Email       string   `json:"email"`
	Role        UserRole `json:"role"`
}
