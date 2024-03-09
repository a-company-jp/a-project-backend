package exception

import "errors"

var (
	ErrorInvalidHeader   = errors.New("INVALID Authorization Header")
	ErrInvalidJWT        = errors.New("INVALID JWT")
	ErrIDAlreadyAssigned = errors.New("ID ALREADY ASSIGNED")
	ErrIDNotAssigned     = errors.New("ID NOT ASSIGNED")
	ErrInvalidRoleLevel  = errors.New("INVALID ROLE LEVEL")
	ErrNotFound          = errors.New("NOT FOUND")
	ErrUnauthorized      = errors.New("UNAUTHORIZED")
	ErrAlreadyExists     = errors.New("ALREADY EXISTS")
)
