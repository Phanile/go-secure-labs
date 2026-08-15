package seed

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunSeed(ctx context.Context, pool *pgxpool.Pool) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	if err := seedCategories(ctx, tx); err != nil {
		return err
	}

	if err := seedProducts(ctx, tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
