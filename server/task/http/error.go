package http

import "errors"

var (
	ErrTaskNotFound       = errors.New("task not found")
	ErrTaskAlreadyExists  = errors.New("task already exists")
	ErrDontHavePermission = errors.New("you don't have permission")
	ErrInvalidFormat      = errors.New("invalid date format, try <yyyy-mm-ddThh:mm:ssZ>")
)
