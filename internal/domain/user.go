package domain

import (
	"fmt"
	"time"

	core_errors "github.com/Phanile/go-secure-labs/internal/errors"
)

type User struct {
	Id           int64
	Email        string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(id int64, email, passwordHash, role string) User {
	return User{
		Id:           id,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func NewUserUninitialized(email, passwordHash, role string) User {
	return User{
		Id:           UninitializedID,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func (u User) Validate() error {
	if u.Email == "" {
		return fmt.Errorf("validate user email: %w", core_errors.ErrInvalidArgument)
	}

	if u.PasswordHash == "" {
		return fmt.Errorf("validate user password hash: %w", core_errors.ErrInvalidArgument)
	}

	if u.Role == "" {
		return fmt.Errorf("validate user role: %w", core_errors.ErrInvalidArgument)
	}

	if u.Id <= 0 {
		return fmt.Errorf("validate user id: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}
