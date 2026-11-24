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
	for _, variantRequest := range req.Variants {
		isActive := true
		if variantRequest.IsActive != nil {
			isActive = *variantRequest.IsActive
		}
		product.Variants = append(product.Variants, domain.ProductVariant{
			Size: variantRequest.Size, Color: variantRequest.Color,
			Stock: variantRequest.Stock, IsActive: isActive,
		})
	}

	if err := s.productRepo.Create(ctx, &product); err != nil {
		return nil, err
	}

	response := toProductResponse(&product)
	return &response, nil
}

func (s *productService) GetProducts(ctx context.Context, param domain.ProductQueryParam) (*domain.ProductPaginationResponse, error) {
	products, totalData, err := s.productRepo.FindAll(ctx, param)
	if err != nil {
		return nil, err
	}

	var productResponses []domain.ProductResponse
	for _, p := range products {
		productResponses = append(productResponses, toProductResponse(&p))
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

	response := toProductResponse(p)
	return &response, nil
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

	response := toProductResponse(p)
	return &response, nil
}

func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.productRepo.Delete(ctx, parsedID)
}

func toProductResponse(product *domain.Product) domain.ProductResponse {
	var category *domain.CategoryResponse
	if product.Category != nil {
		category = &domain.CategoryResponse{
			ID: product.Category.ID, Name: product.Category.Name,
			Slug: product.Category.Slug, Description: product.Category.Description,
			Image: product.Category.Image, CreatedAt: product.Category.CreatedAt,
			UpdatedAt: product.Category.UpdatedAt,
		}
	}

	variants := make([]domain.ProductVariantResponse, 0, len(product.Variants))
	for _, variant := range product.Variants {
		variants = append(variants, domain.ProductVariantResponse{
			ID: variant.ID, ProductID: variant.ProductID, Size: variant.Size,
			Color: variant.Color, Stock: variant.Stock, IsActive: variant.IsActive,
		})
	}

	return domain.ProductResponse{
		ID: product.ID, CategoryID: product.CategoryID, Category: category,
		Variants: variants, Title: product.Title, Price: product.Price,
		OriginalPrice: product.OriginalPrice, Description: product.Description,
		Image: product.Image, RatingRate: product.RatingRate,
		RatingCount: product.RatingCount, CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}
