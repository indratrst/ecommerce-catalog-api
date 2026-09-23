package service

import (
	"context"
	"ecommerce-catalog-api/domain"
	"math"
	"strings"

	"github.com/google/uuid"
)

type productService struct {
	productRepo domain.ProductRepository
}

func NewProductService(productRepo domain.ProductRepository) domain.ProductService {
	return &productService{productRepo: productRepo}
}

func (s *productService) CreateProduct(ctx context.Context, req domain.ProductRequest) (*domain.ProductResponse, error) {
	categoryUUID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return nil, err
	}

	slug := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(req.Name), " ", "-"))
	product := domain.Product{
		CategoryID:  categoryUUID,
		Name:        strings.TrimSpace(req.Name),
		Slug:        slug,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	if err := s.productRepo.Create(ctx, &product); err != nil {
		return nil, err
	}

	return &domain.ProductResponse{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		CreatedAt:   product.CreatedAt,
	}, nil
}

func (s *productService) GetProducts(ctx context.Context, param domain.ProductQueryParam) (*domain.ProductPaginationResponse, error) {
	products, totalData, err := s.productRepo.FindAll(ctx, param)
	if err != nil {
		return nil, err
	}

	var productResponses []domain.ProductResponse
	for _, p := range products {
		productResponses = append(productResponses, domain.ProductResponse{
			ID:          p.ID,
			CategoryID:  p.CategoryID,
			Name:        p.Name,
			Slug:        p.Slug,
			Description: p.Description,
			Price:       p.Price,
			Stock:       p.Stock,
			CreatedAt:   p.CreatedAt,
		})
	}

	limit := param.Limit
	if limit <= 0 {
		limit = 10
	}
	page := param.Page
	if page <= 0 {
		page = 1
	}
	totalPages := int(math.Ceil(float64(totalData) / float64(limit)))

	return &domain.ProductPaginationResponse{
		Data:        productResponses,
		TotalData:   totalData,
		TotalPages:  totalPages,
		CurrentPage: page,
		Limit:       limit,
	}, nil
}

func (s *productService) GetProductByID(ctx context.Context, id string) (*domain.ProductResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	p, err := s.productRepo.FindByID(ctx, parsedID)
	if err != nil {
		return nil, err
	}

	return &domain.ProductResponse{
		ID:          p.ID,
		CategoryID:  p.CategoryID,
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		CreatedAt:   p.CreatedAt,
	}, nil
}
func (s *productService) UpdateProduct(ctx context.Context, id string, req domain.ProductRequest) (*domain.ProductResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	p, err := s.productRepo.FindByID(ctx, parsedID)
	if err != nil {
		return nil, err
	}

	if req.CategoryID != "" {
		catUUID, err := uuid.Parse(req.CategoryID)
		if err != nil {
			return nil, err
		}
		p.CategoryID = catUUID
	}

	p.Name = req.Name
	p.Description = req.Description
	p.Price = req.Price
	p.Stock = req.Stock

	if err := s.productRepo.Update(ctx, p); err != nil {
		return nil, err
	}

	return &domain.ProductResponse{
		ID:          p.ID,
		CategoryID:  p.CategoryID,
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		CreatedAt:   p.CreatedAt,
	}, nil
}

func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.productRepo.Delete(ctx, parsedID)
}
