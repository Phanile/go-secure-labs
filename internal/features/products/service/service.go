package products_service

import (
	"context"

	"github.com/Phanile/go-secure-labs/internal/domain"
)

type ProductsService struct {
	productsRepository ProductsRepository
}

type ProductFilter struct {
	Search string
}

type ProductsRepository interface {
	CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error)
	GetProducts(ctx context.Context, filter ProductFilter) ([]domain.Product, error)
}

func NewProductsService(productsRepository ProductsRepository) *ProductsService {
	return &ProductsService{
		productsRepository: productsRepository,
	}
}
