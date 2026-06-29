package core_errors

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrInvalidArgument    = errors.New("invalid argument")
	ErrConflict           = errors.New("conflict")
	ErrUnauthorized       = errors.New("invalid credentials")
	ErrNoRights           = errors.New("no rights for this action")
	ErrBadGateway         = errors.New("bad gateway")
	ErrServiceUnavailable = errors.New("service unavailable")
)
