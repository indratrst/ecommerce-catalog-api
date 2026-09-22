package service

import (
	"ecommerce-catalog-api/domain"
	"ecommerce-catalog-api/utils"
	"errors"
)

type authService struct {
	userRepo domain.UserRepository
}

func NewAuthService(userRepo domain.UserRepository) domain.AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(req domain.RegisterRequest) (domain.AuthResponse, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return domain.AuthResponse{}, err
	}

	user := domain.User{
		Name:        req.Name,
		Email:       req.Email,
		MobilePhone: req.MobilePhone,
		Password:    hashedPassword,
	}

	if err := s.userRepo.Create(&user); err != nil {
		return domain.AuthResponse{}, errors.New("email sudah terdaftar")
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return domain.AuthResponse{}, err
	}

	return domain.AuthResponse{Token: token, User: user}, nil
}

func (s *authService) Login(req domain.LoginRequest) (domain.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return domain.AuthResponse{}, errors.New("email atau password salah")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return domain.AuthResponse{}, errors.New("email atau password salah")
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return domain.AuthResponse{}, err
	}

	user.Token = token
	s.userRepo.Update(&user)
	return domain.AuthResponse{Token: token, User: user}, nil
}
