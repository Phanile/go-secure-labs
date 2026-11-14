package auth_service

import "context"

func (s *AuthService) Logout(ctx context.Context, userId int64) error {
	err := s.sessionHandler.DeleteSessionByUserId(ctx, userId)
	if err != nil {
		return err
	}

	return nil
}
