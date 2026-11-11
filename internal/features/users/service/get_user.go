package users_service

import (
	"context"

	"github.com/Phanile/go-secure-labs/internal/domain"
)

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
