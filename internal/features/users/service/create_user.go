package users_service

import (
	"context"
	"fmt"

	"github.com/Phanile/go-secure-labs/internal/domain"
	"github.com/Phanile/go-secure-labs/internal/dto"
)

func (s *UserService) CreateUser(ctx context.Context, createUserDTO dto.CreateUserDTO) (domain.User, error) {
	uninitializedUser := domain.NewUserUninitialized(createUserDTO.Email, HashPassword(createUserDTO.Password))
	user, err := s.userRepository.CreateUser(ctx, uninitializedUser)

	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}
