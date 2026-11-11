package domain

import "time"

type Session struct {
	Id         int64
	UserId     int64
	TokenHash  string
	Token      string
	ExpiresAt  time.Time
	CreatedAt  time.Time
	LastUsedAt time.Time
}

func NewSession(id, userId int64, token, tokenHash string, expiresAt, createdAt, lastUsedAt time.Time) Session {
	return Session{
		Id:         id,
		UserId:     userId,
		Token:      token,
		TokenHash:  tokenHash,
		ExpiresAt:  expiresAt,
		CreatedAt:  createdAt,
		LastUsedAt: lastUsedAt,
	}
}

func NewSessionUninitialized(userId int64, token, tokenHash string, expiresAt time.Time) Session {
	return Session{
		Id:         UninitializedID,
		UserId:     userId,
		Token:      token,
		TokenHash:  tokenHash,
		ExpiresAt:  expiresAt,
		CreatedAt:  time.Now(),
		LastUsedAt: time.Now(),
	}
}

func (s Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
