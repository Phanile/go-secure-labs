package seed

import (
	"context"

	core_postgres_pool "github.com/Phanile/go-secure-labs/internal/repository/postgres/pool"
)

func RunSeed(ctx context.Context, pool *core_postgres_pool.ConnectionPool) error {
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
