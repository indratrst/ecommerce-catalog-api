package service

import (
	"context"
	"errors"
	"strings"

	"ecommerce-catalog-api/domain"
	"ecommerce-catalog-api/utils"
)

type authService struct {
	userRepo domain.UserRepository
}

func NewAuthService(userRepo domain.UserRepository) domain.AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error) {
	// 1. Cek ketersediaan email
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email sudah terdaftar")
	}

	// 2. Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("gagal memproses password")
	}

	user := domain.User{
		Name:     strings.TrimSpace(req.Name),
		Email:    strings.ToLower(strings.TrimSpace(req.Email)),
		Password: hashedPassword,
		Role:     domain.RoleUser,
	}

	// 5. Simpan ke DB
	if err := s.userRepo.Create(ctx, &user); err != nil {
		// log.Println("ERROR DB ASLI Saat Create User:", err) // <-- TAMBAHKAN INI UNTUK DEBUG
		return nil, err // Kembalikan error aslinya dulu
	}

	// 6. Generate Token dengan uuid.UUID dan Role
	token, err := utils.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return nil, errors.New("gagal membuat token autentikasi")
	}

	return &domain.AuthResponse{
		Token: token,
		User: domain.UserResponse{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
			Image: user.Image,
		},
	}, nil
}

func (s *authService) Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error) {
	// 1. Cari user berdasarkan email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("email atau password salah")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("email atau password salah")
	}

	// 3. Generate Token
	token, err := utils.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return nil, errors.New("gagal membuat token autentikasi")
	}

	// 4. Extract phone number value dari pointer jika ada

	return &domain.AuthResponse{
		Token: token,
		User: domain.UserResponse{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
			Image: user.Image,
		},
	}, nil
}
