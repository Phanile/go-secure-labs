package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/Phanile/go-secure-labs/internal/domain"
)

func (r *UserRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO appleshop.users (email, password_hash, role)
		VALUES ($1, $2, $3)
		RETURNING id, email, password_hash, role
	`

	row := r.pool.QueryRow(ctx, query, user.Email, user.PasswordHash, user.Role)

	var userModel UserModel
	err := row.Scan(
		&userModel.Id,
		&userModel.Email,
		&userModel.PasswordHash,
		&userModel.Role,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("scan create user: %w", err)
	}

	userDomain := domain.NewUser(
		userModel.Id,
		userModel.Email,
		userModel.PasswordHash,
		userModel.Role,
	)

	return userDomain, nil
}
