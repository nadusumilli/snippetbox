package errors

import (
	"errors"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
	ErrInvalidEmail       = errors.New("models: invalid email address")
	ErrInvalidPassword    = errors.New("models: invalid password")
)
