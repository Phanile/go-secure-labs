package users_service

import (
	"context"
	"fmt"

	"github.com/Phanile/go-secure-labs/internal/domain"
	"github.com/Phanile/go-secure-labs/internal/dto"
)

const DefaultUserRole = "user"

func (s *UserService) CreateUser(ctx context.Context, createUserDTO dto.CreateUserDTO) (domain.User, error) {
	hashedPassword, err := s.passwordHasher.Hash(createUserDTO.Password)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	uninitializedUser := domain.NewUserUninitialized(createUserDTO.Email, hashedPassword, DefaultUserRole)
	user, err := s.userRepository.CreateUser(ctx, uninitializedUser)

	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}
