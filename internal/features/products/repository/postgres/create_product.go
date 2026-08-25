package products_postgres_repository

import (
	"context"
	"fmt"

	"github.com/Phanile/go-secure-labs/internal/domain"
)

func (r *ProductsRepository) CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO appleshop.products (name, description, price, stock_quantity, category_id, image_url) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, description, price, stock_quantity, category_id, image_url
	`

	row := r.pool.QueryRow(ctx, query, product.Name, product.Description, product.Price, product.StockQuantity, product.CategoryId, product.ImageUrl)

	var productModel ProductModel
	err := row.Scan(
		&productModel.Id,
		&productModel.Name,
		&productModel.Description,
		&productModel.Price,
		&productModel.StockQuantity,
		&productModel.CategoryId,
		&productModel.ImageUrl,
	)
	if err != nil {
		return domain.Product{}, fmt.Errorf("scan create product: %w", err)
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

	return productDomain, nil
}
