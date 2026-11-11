package auth_service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Phanile/go-secure-labs/internal/domain"
	"github.com/Phanile/go-secure-labs/internal/dto"
)

const SessionLiveTime = 30 * 24 * time.Hour

func (s *AuthService) Register(ctx context.Context, createUserDTO dto.CreateUserDTO) (domain.User, domain.Session, error) {
	user, err := s.userHandler.CreateUser(ctx, createUserDTO)
	if err != nil {
		return domain.User{}, domain.Session{}, err
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
