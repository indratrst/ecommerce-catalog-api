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
	slug := categorySlug(req.Name, req.Slug)
	category := domain.Category{
		Name:        strings.TrimSpace(req.Name),
		Slug:        slug,
		Description: req.Description,
		Image:       req.Image,
	}

	if err := s.categoryRepo.Create(ctx, &category); err != nil {
		return nil, err
	}

	return &domain.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		Image:       category.Image,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
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
			ID:          c.ID,
			Name:        c.Name,
			Slug:        c.Slug,
			Description: c.Description,
			Image:       c.Image,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
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
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		Image:       category.Image,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, id string, req domain.CategoryRequest) (*domain.CategoryResponse, error) {
	categoryID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	category, err := s.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		category.Name = strings.TrimSpace(req.Name)
	}
	if req.Slug != "" || req.Name != "" {
		category.Slug = categorySlug(category.Name, req.Slug)
	}
	category.Description = req.Description
	category.Image = req.Image

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	return &domain.CategoryResponse{
		ID: category.ID, Name: category.Name, Slug: category.Slug,
		Description: category.Description, Image: category.Image,
		CreatedAt: category.CreatedAt, UpdatedAt: category.UpdatedAt,
	}, nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, id string) error {
	categoryID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.categoryRepo.Delete(ctx, categoryID)
}

func categorySlug(name, slug string) string {
	if slug = strings.TrimSpace(slug); slug != "" {
		return strings.ToLower(strings.ReplaceAll(slug, " ", "-"))
	}
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), " ", "-"))
}
