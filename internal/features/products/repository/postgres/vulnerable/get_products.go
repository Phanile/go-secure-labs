package products_postgres_vulnerable_repository

import (
	"context"
	"fmt"

	"github.com/Phanile/go-secure-labs/internal/domain"
	products_postgres_repository "github.com/Phanile/go-secure-labs/internal/features/products/repository/postgres"
	products_service "github.com/Phanile/go-secure-labs/internal/features/products/service"
)

func (r *ProductsVulnerableRepository) GetProducts(ctx context.Context, filter products_service.ProductFilter) ([]domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var query string
	if filter.Search != "" {
		query = fmt.Sprintf(
			`SELECT id, name, description, price, stock_quantity, category_id, image_url 
			FROM appleshop.products
			WHERE name ILIKE '%%%s%%'`, filter.Search)
	} else {
		query = `
			SELECT id, name, description, price, stock_quantity, category_id, image_url 
			FROM appleshop.products
		`
	}

	rows, err := r.pool.Query(ctx, query)

	if err != nil {
		return []domain.Product{}, fmt.Errorf("getproducts: %w", err)
	}
	defer rows.Close()

	products := make([]domain.Product, 0)
	for rows.Next() {
		var productModel products_postgres_repository.ProductModel
		err := rows.Scan(
			&productModel.Id,
			&productModel.Name,
			&productModel.Description,
			&productModel.Price,
			&productModel.StockQuantity,
			&productModel.CategoryId,
			&productModel.ImageUrl,
		)

		if err != nil {
			return []domain.Product{}, fmt.Errorf("getproducts scan row: %w", err)
		}

		productDomain := domain.NewProduct(
			productModel.Id,
			productModel.Name,
			productModel.Description,
			productModel.Price,
			productModel.StockQuantity,
			productModel.CategoryId,
			productModel.ImageUrl,
		)

		products = append(products, productDomain)
	}

	if rows.Err() != nil {
		return []domain.Product{}, fmt.Errorf("getproducts scan rows: %w", err)
	}

	return products, nil
}
