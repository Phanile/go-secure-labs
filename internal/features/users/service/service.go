package users_service

import (
	"context"

	"github.com/Phanile/go-secure-labs/internal/domain"
)

type UserService struct {
	userRepository UserRepository
	passwordHasher PasswordHasher
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password, hash string) (bool, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
}

func NewUserService(userRepository UserRepository, passHasher PasswordHasher) *UserService {
	return &UserService{
		userRepository: userRepository,
		passwordHasher: passHasher,
	}
}
