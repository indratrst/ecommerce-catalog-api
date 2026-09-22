package domain

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name     string    `json:"name" gorm:"not null"`
	Products []Product `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
}

type CategoryRequest struct {
	Name string `json:"name"`
}

type CategoryRepository interface {
	Create(category *Category) error
	FindAll() ([]Category, error)
	FindByID(id uint) (Category, error)
}

type CategoryService interface {
	CreateCategory(req CategoryRequest) (Category, error)
	GetAllCategories() ([]Category, error)
}
