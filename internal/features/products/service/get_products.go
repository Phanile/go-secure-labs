package products_service

import (
	"context"
	"fmt"

	"github.com/Phanile/go-secure-labs/internal/domain"
)

func (s *ProductsService) GetProducts(ctx context.Context) ([]domain.Product, error) {
	products, err := s.productsRepository.GetProducts(ctx)

	if err != nil {
		return []domain.Product{}, fmt.Errorf("get products: %w", err)
	}

	return products, nil
}
