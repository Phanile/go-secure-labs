package auth_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Phanile/go-secure-labs/internal/domain"
	internal_errors "github.com/Phanile/go-secure-labs/internal/errors"
	"github.com/jackc/pgx/v5"
)

const EmptyToken = ""

func (r *AuthRepository) GetSessionByTokenHash(ctx context.Context, tokenHash string) (domain.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var query string

	query = `SELECT id, user_id, token_hash, expires_at, created_at, last_used_at 
			FROM appleshop.sessions
			WHERE token_hash = $1`

	var sessionModel SessionModel

	err := r.pool.QueryRow(ctx, query, tokenHash).Scan(
		&sessionModel.Id,
		&sessionModel.UserId,
		&sessionModel.TokenHash,
		&sessionModel.ExpiresAt,
		&sessionModel.CreatedAt,
		&sessionModel.LastUsedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Session{}, fmt.Errorf("session with token hash=%s not found: %w", tokenHash, internal_errors.ErrNotFound)
		}
		return domain.Session{}, err
	}

	sessionDomain := domain.NewSession(
		sessionModel.Id,
		sessionModel.UserId,
		EmptyToken,
		sessionModel.TokenHash,
		sessionModel.ExpiresAt,
		sessionModel.CreatedAt,
		sessionModel.LastUsedAt,
	)

	return sessionDomain, nil
}
