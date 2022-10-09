package localstore

import "errors"

var (
	ErrTaskNotFound       = errors.New("task not found")
	ErrTaskAlreadyExists  = errors.New("task already exists")
	ErrDontHavePermission = errors.New("you don't have permission")
	ErrInvalidStatus      = errors.New("invalid status, valid: <in progress>, <todo>, <done>")
)
