package seed

import (
	"context"

	"github.com/Phanile/go-secure-labs/internal/domain"
	"github.com/jackc/pgx/v5"
)

func seedProducts(ctx context.Context, tx pgx.Tx) error {
	products := []domain.Product{
		domain.NewProduct(1, "Iphone 15", "Good Phone.", 52000, 12, 1, "/images/iphone15.jpg"),
		domain.NewProduct(2, "Iphone 16", "Good Phone.", 63000, 20, 1, "/images/iphone16.jpg"),
		domain.NewProduct(3, "Iphone 17", "Good Phone.", 83000, 32, 1, "/images/iphone17.jpg"),
	}

	batch := &pgx.Batch{}

	for _, product := range products {
		batch.Queue(
			"INSERT INTO appleshop.products (id, name, description, price, stock_quantity, category_id, image_url) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (id) DO NOTHING",
			product.Id,
			product.Name,
			product.Description,
			product.Price,
			product.StockQuantity,
			product.CategoryId,
			product.ImageUrl,
		)
	}

	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	for range products {
		if _, err := results.Exec(); err != nil {
			return err
		}
	}

	return nil
}
