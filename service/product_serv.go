package service

import (
	"ecommerce-catalog-api/domain"
	"math"
)

type productService struct {
	productRepo domain.ProductRepository
}

func NewProductService(productRepo domain.ProductRepository) domain.ProductService {
	return &productService{productRepo: productRepo}
}

func (s *productService) CreateProduct(req domain.ProductRequest) (domain.Product, error) {
	product := domain.Product{
		Name:       req.Name,
		Price:      req.Price,
		Stock:      req.Stock,
		CategoryID: req.CategoryID,
	}

	err := s.productRepo.Create(&product)
	return product, err
}

func (s *productService) GetAllProducts(param domain.ProductQueryParam) (domain.ProductPaginationResponse, error) {
	// Set default jika query param tidak diisi
	if param.Page <= 0 {
		param.Page = 1
	}
	if param.Limit <= 0 {
		param.Limit = 10 // Default 10 item per halaman
	}

	products, totalData, err := s.productRepo.FindAll(param)
	if err != nil {
		return domain.ProductPaginationResponse{}, err
	}

	// Hitung total halaman (contoh: 12 data / limit 10 = 2 halaman)
	totalPages := int(math.Ceil(float64(totalData) / float64(param.Limit)))

	response := domain.ProductPaginationResponse{
		Data: products,
		Meta: domain.MetaPagination{
			CurrentPage: param.Page,
			TotalPages:  totalPages,
			Limit:       param.Limit,
			TotalData:   totalData,
		},
	}

	return response, nil
}

func (s *productService) GetProductByID(id string) (domain.Product, error) {
	return s.productRepo.FindByID(id)
}

func (s *productService) UpdateProduct(id string, req domain.ProductRequest) (domain.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return product, err
	}

	product.Name = req.Name
	product.Price = req.Price
	product.Stock = req.Stock

	err = s.productRepo.Update(&product)
	return product, err
}

func (s *productService) DeleteProduct(id string) error {
	return s.productRepo.Delete(id)
}
