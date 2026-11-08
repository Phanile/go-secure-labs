package users_postgres_repository

import core_postgres_pool "github.com/Phanile/go-secure-labs/internal/repository/postgres/pool"

type UserRepository struct {
	pool core_postgres_pool.Pool
}

func NewUserRepository(p core_postgres_pool.Pool) *UserRepository {
	return &UserRepository{
		pool: p,
	}
}
