package auth_service

import (
	"context"

	"github.com/Phanile/go-secure-labs/internal/domain"
	"github.com/Phanile/go-secure-labs/internal/dto"
)

func (s *AuthService) Register(ctx context.Context, createUserDTO dto.CreateUserDTO) (domain.User, domain.Session, error) {
	user, err := s.userCreator.CreateUser(ctx, createUserDTO)
	if err != nil {
		return domain.User{}, domain.Session{}, err
	}

	session, err := s.sessionCreator.CreateSession(ctx, user.Id)
	if err != nil {
		return domain.User{}, domain.Session{}, err
	}

	return user, session, nil
}
