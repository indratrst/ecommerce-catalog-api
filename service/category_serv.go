package service

import (
	"context"
	"ecommerce-catalog-api/domain"
	"strings"

	"github.com/google/uuid"
)

type categoryService struct {
	categoryRepo domain.CategoryRepository
}

func NewCategoryService(categoryRepo domain.CategoryRepository) domain.CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func (s *categoryService) CreateCategory(ctx context.Context, req domain.CategoryRequest) (*domain.CategoryResponse, error) {
	slug := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(req.Name), " ", "-"))
	category := domain.Category{
		Name: strings.TrimSpace(req.Name),
		Slug: slug,
	}

	if err := s.categoryRepo.Create(ctx, &category); err != nil {
		return nil, err
	}

	return &domain.CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		CreatedAt: category.CreatedAt,
	}, nil
}

func (s *categoryService) GetCategories(ctx context.Context) ([]domain.CategoryResponse, error) {
	categories, err := s.categoryRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var response []domain.CategoryResponse
	for _, c := range categories {
		response = append(response, domain.CategoryResponse{
			ID:        c.ID,
			Name:      c.Name,
			Slug:      c.Slug,
			CreatedAt: c.CreatedAt,
		})
	}
	return response, nil
}

func (s *categoryService) GetCategoryByID(ctx context.Context, id string) (*domain.CategoryResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	category, err := s.categoryRepo.FindByID(ctx, parsedID)
	if err != nil {
		return nil, err
	}

	return &domain.CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		CreatedAt: category.CreatedAt,
	}, nil
}
