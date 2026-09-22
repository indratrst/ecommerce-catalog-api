package repository

import (
	"my-go-api/domain"
	"strings"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *domain.Product) error {
	// Lakukan Create terlebih dahulu
	err := r.db.Create(product).Error
	if err != nil {
		return err
	}

	// Ambil kembali data product beserta relasi Category yang baru dibuat
	return r.db.Preload("Category").First(product, product.ID).Error
}

func (r *productRepository) FindAll(param domain.ProductQueryParam) ([]domain.Product, int64, error) {
	var products []domain.Product
	var totalData int64

	// Mulai query dasar
	query := r.db.Model(&domain.Product{})

	// 1. Filter Search berdasarkan nama produk (Case-insensitive)
	if param.Search != "" {
		searchTerm := "%" + strings.ToLower(param.Search) + "%"
		query = query.Where("LOWER(name) LIKE ?", searchTerm)
	}

	// 2. Filter opsional berdasarkan category_id
	if param.CategoryID > 0 {
		query = query.Where("category_id = ?", param.CategoryID)
	}

	// Hitung total data sebelum diterapkan offset/limit pagination
	if err := query.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	// 3. Terapkan Pagination (Offset & Limit)
	offset := (param.Page - 1) * param.Limit
	err := query.Preload("Category").
		Limit(param.Limit).
		Offset(offset).
		Order("id DESC"). // Urutkan dari produk terbaru
		Find(&products).Error

	return products, totalData, err
}

func (r *productRepository) FindByID(id string) (domain.Product, error) {
	var product domain.Product
	err := r.db.Preload("Category").First(&product, id).Error
	return product, err
}

func (r *productRepository) Update(product *domain.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id string) error {
	return r.db.Delete(&domain.Product{}, id).Error
}
