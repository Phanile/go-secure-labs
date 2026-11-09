package domain

import "time"

type Session struct {
	Id         int64
	UserId     int
	TokenHash  string
	ExpiresAt  time.Time
	CreatedAt  time.Time
	LastUsedAt time.Time
}
