package application

import "errors"

var (
	// application errors.
	ErrUserAlreadyExists  = errors.New("user with this login already exists")
	ErrInvalidCredentials = errors.New("invalid login or password")
	ErrUnvailible         = errors.New("service unvailible")

	// infra errors.
	ErrNotFound = errors.New("not found")
)
