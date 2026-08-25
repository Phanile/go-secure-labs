package products_postgres_repository

import core_postgres_pool "github.com/Phanile/go-secure-labs/internal/repository/postgres/pool"

type ProductsRepository struct {
	pool core_postgres_pool.Pool
}

func NewProductsRepository(p core_postgres_pool.Pool) *ProductsRepository {
	return &ProductsRepository{}
}
