package psql

import "errors"

var (
	ErrWrongLogin         = errors.New("invalid login, try another one")
	ErrWrongPassword      = errors.New("invalid password, try another one")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidAccessToken = errors.New("invalid access token")
)
