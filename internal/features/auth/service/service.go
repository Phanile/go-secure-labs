package auth_service

import (
	"context"

	"github.com/Phanile/go-secure-labs/internal/domain"
	"github.com/Phanile/go-secure-labs/internal/dto"
)

type AuthService struct {
	userCreator    UserCreator
	sessionCreator SessionCreator
}

type UserCreator interface {
	CreateUser(ctx context.Context, createUserDTO dto.CreateUserDTO) (domain.User, error)
}

type SessionCreator interface {
	CreateSession(ctx context.Context, userId int64) (domain.Session, error)
}

func NewAuthService(userCreator UserCreator, sessionCreator SessionCreator) *AuthService {
	return &AuthService{
		userCreator:    userCreator,
		sessionCreator: sessionCreator,
	}
}
