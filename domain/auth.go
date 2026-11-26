package domain

import (
	"context"
	"time"
)

// 1. Interface AuthService (Contract untuk Business Logic Auth)
type AuthService interface {
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindAll(ctx context.Context) ([]User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type AdminUserRequest struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password string   `json:"password,omitempty"`
	Role     UserRole `json:"role"`
	Image    *string  `json:"image"`
}

type AdminUserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      UserRole  `json:"role"`
	Image     *string   `json:"image,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

type AdminUserService interface {
	GetUsers(ctx context.Context) ([]AdminUserResponse, error)
	CreateUser(ctx context.Context, req AdminUserRequest) (*AdminUserResponse, error)
	UpdateUser(ctx context.Context, id string, req AdminUserRequest) (*AdminUserResponse, error)
	DeleteUser(ctx context.Context, id string) error
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
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
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Role  UserRole `json:"role"`
	Image *string  `json:"image,omitempty"`
}
