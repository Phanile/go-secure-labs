package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Phanile/go-secure-labs/internal/domain"
	internal_errors "github.com/Phanile/go-secure-labs/internal/errors"
	"github.com/jackc/pgx/v5"
)

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var query string

	query =
		`SELECT id, email, password_hash, role 
		FROM appleshop.users
		WHERE email = $1`

	var userModel UserModel

	err := r.pool.QueryRow(ctx, query, email).Scan(
		&userModel.Id,
		&userModel.Email,
		&userModel.PasswordHash,
		&userModel.Role,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with email=%s not found: %w", email, internal_errors.ErrNotFound)
		}
		return domain.User{}, err
	}

	userDomain := domain.NewUser(
		userModel.Id,
		userModel.Email,
		userModel.PasswordHash,
		userModel.Role,
	)

	return userDomain, nil
}
