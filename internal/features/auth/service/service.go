package auth_service

import (
	"context"

	"github.com/Phanile/go-secure-labs/internal/domain"
	"github.com/Phanile/go-secure-labs/internal/dto"
)

type AuthService struct {
	userHandler    UserHandler
	sessionHandler SessionHandler
	tokenHasher    TokenHasher
	passwordHasher PasswordHasher
}

type UserHandler interface {
	CreateUser(ctx context.Context, createUserDTO dto.CreateUserDTO) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
}

type SessionHandler interface {
	CreateSession(ctx context.Context, session domain.Session) (domain.Session, error)
}

type TokenHasher interface {
	Hash(token string) (string, error)
	Compare(token, hash string) (bool, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password, hash string) (bool, error)
}

func NewAuthService(userHandler UserHandler, sessionHandler SessionHandler, tokenHasher TokenHasher, passwordHasher PasswordHasher) *AuthService {
	return &AuthService{
		userHandler:    userHandler,
		sessionHandler: sessionHandler,
		tokenHasher:    tokenHasher,
		passwordHasher: passwordHasher,
	}
}
