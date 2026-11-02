package service

import (
	"context"
	"ecommerce-catalog-api/domain"
	"math"

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

	product := domain.Product{
		CategoryID:    categoryUUID,
		Title:         req.Title,
		Price:         req.Price,
		OriginalPrice: req.OriginalPrice,
		Description:   req.Description,
		Image:         req.Image,
	}

	if err := s.productRepo.Create(ctx, &product); err != nil {
		return nil, err
	}

	return &domain.ProductResponse{
		ID:            product.ID,
		CategoryID:    product.CategoryID,
		Title:         product.Title,
		Price:         product.Price,
		OriginalPrice: product.OriginalPrice,
		Description:   product.Description,
		Image:         product.Image,
		RatingRate:    product.RatingRate,
		RatingCount:   product.RatingCount,
		CreatedAt:     product.CreatedAt,
		UpdatedAt:     product.UpdatedAt,
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
			ID:            p.ID,
			CategoryID:    p.CategoryID,
			Title:         p.Title,
			Price:         p.Price,
			OriginalPrice: p.OriginalPrice,
			Description:   p.Description,
			Image:         p.Image,
			RatingRate:    p.RatingRate,
			RatingCount:   p.RatingCount,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
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
		ID:            p.ID,
		CategoryID:    p.CategoryID,
		Title:         p.Title,
		Price:         p.Price,
		OriginalPrice: p.OriginalPrice,
		Description:   p.Description,
		Image:         p.Image,
		RatingRate:    p.RatingRate,
		RatingCount:   p.RatingCount,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
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

	p.Title = req.Title
	p.Description = req.Description
	p.Price = req.Price
	p.OriginalPrice = req.OriginalPrice
	p.Image = req.Image

	if err := s.productRepo.Update(ctx, p); err != nil {
		return nil, err
	}

	return &domain.ProductResponse{
		ID:            p.ID,
		CategoryID:    p.CategoryID,
		Title:         p.Title,
		Price:         p.Price,
		OriginalPrice: p.OriginalPrice,
		Description:   p.Description,
		Image:         p.Image,
		RatingRate:    p.RatingRate,
		RatingCount:   p.RatingCount,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}, nil
}

func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.productRepo.Delete(ctx, parsedID)
}
