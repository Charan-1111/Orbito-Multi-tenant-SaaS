package constants

import "errors"

var (
	ErrMissingJwtSecret    = errors.New("JWT secret is missing")
	ErrInvalidToken        = errors.New("Invalid token")
	ErrInvalidRefreshToken = errors.New("Invalid refresh token")
	ErrUserNotFound        = errors.New("User not found")
	ErrLogoutFailed        = errors.New("Logout failed")
)
