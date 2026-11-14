package auth_postgres_repository

import (
	"context"
	"fmt"
)

func (r *AuthRepository) DeleteSessionByUserId(ctx context.Context, userId int64) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `DELETE FROM appleshop.sessions WHERE user_id = $1`

	if _, err := r.pool.Exec(ctx, query, userId); err != nil {
		return fmt.Errorf("delete sessions for user %d: %w", userId, err)
	}

	return nil
}
