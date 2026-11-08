package errors

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrInvalidArgument  = errors.New("invalid argument")
	ErrConflict         = errors.New("conflict")
	ErrPasswordMismatch = errors.New("password mismatch")
)
