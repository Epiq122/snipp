package models

import "errors"

var (
	ErrNoRecord           = errors.New("models: no matching records found")
	ErrInvalidCredentials = errors.New("models: invalid credentials provided")
	ErrDuplicateEmail     = errors.New("models: duplicate email provided")
)
