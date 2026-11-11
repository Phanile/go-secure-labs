package auth_postgres_repository

import (
	"context"
	"fmt"

	"github.com/Phanile/go-secure-labs/internal/domain"
)

func (r *AuthRepository) CreateSession(ctx context.Context, session domain.Session) (domain.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO appleshop.sessions (user_id, token_hash, expires_at, created_at, last_used_at) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, token_hash, expires_at, created_at, last_used_at
	`

	row := r.pool.QueryRow(ctx, query, session.UserId, session.TokenHash, session.ExpiresAt, session.CreatedAt, session.LastUsedAt)

	var sessionModel SessionModel
	err := row.Scan(
		&sessionModel.Id,
		&sessionModel.UserId,
		&sessionModel.TokenHash,
		&sessionModel.ExpiresAt,
		&sessionModel.CreatedAt,
		&sessionModel.LastUsedAt,
	)
	if err != nil {
		return domain.Session{}, fmt.Errorf("scan create session: %w", err)
	}

	sessionDomain := domain.NewSession(
		sessionModel.Id,
		sessionModel.UserId,
		session.Token,
		sessionModel.TokenHash,
		sessionModel.ExpiresAt,
		sessionModel.CreatedAt,
		sessionModel.LastUsedAt,
	)

	return sessionDomain, nil
}
