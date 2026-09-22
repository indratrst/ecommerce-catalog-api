package repository

import (
	"my-go-api/domain"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) domain.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *domain.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) FindAll() ([]domain.Category, error) {
	var categories []domain.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) FindByID(id uint) (domain.Category, error) {
	var category domain.Category
	err := r.db.First(&category, id).Error
	return category, err
}
