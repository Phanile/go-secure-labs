package seed

import (
	"context"

	"github.com/Phanile/go-secure-labs/internal/domain"
	"github.com/jackc/pgx/v5"
)

func seedCategories(ctx context.Context, tx pgx.Tx) error {
	categories := []domain.Category{
		domain.NewCategory(1, "Iphone", "Iphones"),
		domain.NewCategory(2, "MacBook", "MacBooks"),
		domain.NewCategory(3, "AirPods", "AirPods"),
	}

	batch := &pgx.Batch{}

	for _, category := range categories {
		batch.Queue(
			"INSERT INTO appleshop.categories (id, name, description) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING",
			category.Id,
			category.Name,
			category.Description,
		)
	}

	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	for range categories {
		if _, err := results.Exec(); err != nil {
			return err
		}
	}

	return nil
}
