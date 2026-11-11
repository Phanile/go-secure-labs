package auth_postgres_repository

import "time"

type SessionModel struct {
	Id         int64
	UserId     int64
	Token      string
	TokenHash  string
	ExpiresAt  time.Time
	CreatedAt  time.Time
	LastUsedAt time.Time
}
