package service

import (
	"ecommerce-catalog-api/domain"
)

type categoryService struct {
	categoryRepo domain.CategoryRepository
}

func NewCategoryService(categoryRepo domain.CategoryRepository) domain.CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func (s *categoryService) CreateCategory(req domain.CategoryRequest) (domain.Category, error) {
	category := domain.Category{Name: req.Name}
	err := s.categoryRepo.Create(&category)
	return category, err
}

func (s *categoryService) GetAllCategories() ([]domain.Category, error) {
	return s.categoryRepo.FindAll()
}
