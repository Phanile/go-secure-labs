package auth_service

import (
	"context"
	"time"

	"github.com/Phanile/go-secure-labs/internal/domain"
	"github.com/Phanile/go-secure-labs/internal/dto"
	"github.com/Phanile/go-secure-labs/internal/errors"
	"github.com/google/uuid"
)

func (s *AuthService) Login(ctx context.Context, loginUserDTO dto.LoginUserDTO) (domain.User, domain.Session, error) {
	user, err := s.userHandler.GetUserByEmail(ctx, loginUserDTO.Email)
	if err != nil {
		return domain.User{}, domain.Session{}, err
	}

	ok, err := s.passwordHasher.Compare(loginUserDTO.Password, user.PasswordHash)
	if err != nil {
		return domain.User{}, domain.Session{}, err
	}
	if !ok {
		return domain.User{}, domain.Session{}, errors.ErrPasswordMismatch
	}

	sessionToken := uuid.New().String()

	sessionTokenHash, err := s.tokenHasher.Hash(sessionToken)
	if err != nil {
		return domain.User{}, domain.Session{}, err
	}

	uninitializedSession := domain.NewSessionUninitialized(user.Id, sessionToken, sessionTokenHash, time.Now().Add(SessionLiveTime))

	session, err := s.sessionHandler.CreateSession(ctx, uninitializedSession)
	if err != nil {
		return domain.User{}, domain.Session{}, err
	}

	return user, session, nil
}
