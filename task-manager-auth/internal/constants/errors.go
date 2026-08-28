package constants

import "errors"

var (
	ErrMissingJwtSecret = errors.New("JWT secret is missing")
	ErrInvalidToken     = errors.New("Invalid token")
)
