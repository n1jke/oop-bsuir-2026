package application

import "errors"

var (
	// application errors.
	ErrUserAlreadyExists  = errors.New("user with this login already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrBookNotFound       = errors.New("book not found")
	ErrInvalidCredentials = errors.New("invalid login or password")
	ErrInvalidParams      = errors.New("invalid params")
	ErrUnvailible         = errors.New("service unvailible")

	// infra errors.
	ErrNotFound     = errors.New("not found")
	ErrAlreadyExist = errors.New("already exist")
)
