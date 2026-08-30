package products_service

import (
	"context"
	"fmt"

	"github.com/Phanile/go-secure-labs/internal/domain"
)

func (s *ProductsService) GetProducts(ctx context.Context, filter ProductFilter) ([]domain.Product, error) {
	products, err := s.productsRepository.GetProducts(ctx, filter)

	if err != nil {
		return []domain.Product{}, fmt.Errorf("get products: %w", err)
	}

	return products, nil
}
