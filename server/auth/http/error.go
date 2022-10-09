package http

import "errors"

var (
	ErrNoLogin    = errors.New("login must be not null")
	ErrNoPassword = errors.New("password must be not null")
)
