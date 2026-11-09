package auth_postgres_repository

import core_postgres_pool "github.com/Phanile/go-secure-labs/internal/repository/postgres/pool"

type AuthRepository struct {
	pool core_postgres_pool.Pool
}

func NewAuthRepository(p core_postgres_pool.Pool) *AuthRepository {
	return &AuthRepository{
		pool: p,
	}
}
