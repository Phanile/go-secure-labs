package products_postgres_vulnerable_repository

import core_postgres_pool "github.com/Phanile/go-secure-labs/internal/repository/postgres/pool"

type ProductsVulnerableRepository struct {
	pool core_postgres_pool.Pool
}

func NewProductsVulnerableRepository(p core_postgres_pool.Pool) *ProductsVulnerableRepository {
	return &ProductsVulnerableRepository{
		pool: p,
	}
}
