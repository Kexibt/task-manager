package localstore

import "errors"

var (
	ErrTaskNotFound       = errors.New("task not found")
	ErrTaskAlreadyExists  = errors.New("task already exists")
	ErrDontHavePermission = errors.New("you don't have permission")
	ErrIDRequired         = errors.New("ID must be not null")
	ErrTitleRequired      = errors.New("title must be not null")
	ErrStatusRequired     = errors.New("status must be not null")
	ErrInvalidStatus      = errors.New("invalid status, valid: <in progress>, <todo>, <done>")
)
