package auth_service

import (
	"context"

	"github.com/Phanile/go-secure-labs/internal/domain"
)

func (h *AuthService) Me(ctx context.Context, userId int64) (domain.User, error) {
	return h.userHandler.GetUserById(ctx, userId)
}
