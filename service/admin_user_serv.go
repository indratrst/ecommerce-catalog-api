package service

import (
	"context"
	"errors"
	"strings"

	"ecommerce-catalog-api/domain"
	"ecommerce-catalog-api/utils"
)

type adminUserService struct {
	userRepo domain.UserRepository
}

func NewAdminUserService(userRepo domain.UserRepository) domain.AdminUserService {
	return &adminUserService{userRepo: userRepo}
}

func (s *adminUserService) GetUsers(ctx context.Context) ([]domain.AdminUserResponse, error) {
	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]domain.AdminUserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, *adminUserResponse(&user))
	}
	return responses, nil
}

func (s *adminUserService) CreateUser(ctx context.Context, req domain.AdminUserRequest) (*domain.AdminUserResponse, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Name == "" || req.Email == "" || len(req.Password) < 8 {
		return nil, errors.New("nama, email, dan password minimal 8 karakter wajib diisi")
	}
	if !validUserRole(req.Role) {
		return nil, errors.New("role pengguna tidak valid")
	}

	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email sudah terdaftar")
	}

	password, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	user := domain.User{
		Name: req.Name, Email: req.Email, Password: password, Role: req.Role,
	}
	if err := s.userRepo.Create(ctx, &user); err != nil {
		return nil, err
	}
	return adminUserResponse(&user), nil
}

func (s *adminUserService) UpdateUser(ctx context.Context, id string, req domain.AdminUserRequest) (*domain.AdminUserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Name == "" || req.Email == "" {
		return nil, errors.New("nama dan email wajib diisi")
	}
	if !validUserRole(req.Role) {
		return nil, errors.New("role pengguna tidak valid")
	}
	if req.Email != strings.ToLower(user.Email) {
		exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("email sudah terdaftar")
		}
	}
	if req.Password != "" {
		if len(req.Password) < 8 {
			return nil, errors.New("password minimal 8 karakter")
		}
		user.Password, err = utils.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Role = req.Role
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return adminUserResponse(user), nil
}

func (s *adminUserService) DeleteUser(ctx context.Context, id string) error {
	if _, err := s.userRepo.FindByID(ctx, id); err != nil {
		return err
	}
	return s.userRepo.Delete(ctx, id)
}

func adminUserResponse(user *domain.User) *domain.AdminUserResponse {
	return &domain.AdminUserResponse{
		ID: user.ID.String(), Name: user.Name, Email: user.Email,
		Role: user.Role, Image: user.Image, CreatedAt: user.CreatedAt,
	}
}

func validUserRole(role domain.UserRole) bool {
	return role == domain.RoleSuperuser || role == domain.RoleAdmin || role == domain.RoleUser
}
