package auth_service

import (
	"context"

	internal_errors "github.com/Phanile/go-secure-labs/internal/errors"
)

func (s *AuthService) VerifySession(ctx context.Context, sessionToken string) (int64, error) {
	sessionTokenHash, err := s.tokenHasher.Hash(sessionToken)
	if err != nil {
		return 0, err
	}

	session, err := s.sessionHandler.GetSessionByTokenHash(ctx, sessionTokenHash)
	if err != nil {
		return 0, err
	}

	if session.IsExpired() {
		return 0, internal_errors.ErrTokenExpired
	}

	return session.UserId, nil
}
