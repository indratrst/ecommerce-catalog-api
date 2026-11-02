package service

import (
	"context"
	"ecommerce-catalog-api/domain"

	"github.com/google/uuid"
)

type productVariantService struct {
	productVariantRepo domain.ProductVariantRepository
}

func NewProductVariantService(productVariantRepo domain.ProductVariantRepository) domain.ProductVariantService {
	return &productVariantService{productVariantRepo: productVariantRepo}
}

func (s *productVariantService) CreateProductVariant(ctx context.Context, productID string, req domain.ProductVariantRequest) (*domain.ProductVariantResponse, error) {
	parsedProductID, err := uuid.Parse(productID)
	if err != nil {
		return nil, err
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	variant := domain.ProductVariant{
		ProductID: parsedProductID,
		Size:      req.Size,
		Color:     req.Color,
		Stock:     req.Stock,
		IsActive:  isActive,
	}

	if err := s.productVariantRepo.Create(ctx, &variant); err != nil {
		return nil, err
	}

	return &domain.ProductVariantResponse{
		ID:        variant.ID,
		ProductID: variant.ProductID,
		Size:      variant.Size,
		Color:     variant.Color,
		Stock:     variant.Stock,
		IsActive:  variant.IsActive,
	}, nil
}

func (s *productVariantService) GetProductVariantsByProductID(ctx context.Context, productID string) ([]domain.ProductVariantResponse, error) {
	parsedProductID, err := uuid.Parse(productID)
	if err != nil {
		return nil, err
	}

	variants, err := s.productVariantRepo.FindByProductID(ctx, parsedProductID)
	if err != nil {
		return nil, err
	}

	resp := make([]domain.ProductVariantResponse, 0, len(variants))
	for _, v := range variants {
		resp = append(resp, domain.ProductVariantResponse{
			ID:        v.ID,
			ProductID: v.ProductID,
			Size:      v.Size,
			Color:     v.Color,
			Stock:     v.Stock,
			IsActive:  v.IsActive,
		})
	}
	return resp, nil
}

func (s *productVariantService) GetProductVariantByID(ctx context.Context, id string) (*domain.ProductVariantResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	variant, err := s.productVariantRepo.FindByID(ctx, parsedID)
	if err != nil {
		return nil, err
	}

	return &domain.ProductVariantResponse{
		ID:        variant.ID,
		ProductID: variant.ProductID,
		Size:      variant.Size,
		Color:     variant.Color,
		Stock:     variant.Stock,
		IsActive:  variant.IsActive,
	}, nil
}

func (s *productVariantService) UpdateProductVariant(ctx context.Context, id string, req domain.ProductVariantRequest) (*domain.ProductVariantResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	variant, err := s.productVariantRepo.FindByID(ctx, parsedID)
	if err != nil {
		return nil, err
	}

	if req.Size != "" {
		variant.Size = req.Size
	}
	if req.Color != nil {
		variant.Color = req.Color
	}
	if req.Stock >= 0 {
		variant.Stock = req.Stock
	}
	if req.IsActive != nil {
		variant.IsActive = *req.IsActive
	}

	if err := s.productVariantRepo.Update(ctx, variant); err != nil {
		return nil, err
	}

	return &domain.ProductVariantResponse{
		ID:        variant.ID,
		ProductID: variant.ProductID,
		Size:      variant.Size,
		Color:     variant.Color,
		Stock:     variant.Stock,
		IsActive:  variant.IsActive,
	}, nil
}

func (s *productVariantService) DeleteProductVariant(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.productVariantRepo.Delete(ctx, parsedID)
}
