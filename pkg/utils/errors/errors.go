package errors

import (
	"errors"
)

// Custom errors
var (
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidCredentials    = errors.New("invalid email or password")
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrInvalidToken          = errors.New("invalid or expired token")
	ErrSubscriptionExpired   = errors.New("subscription has expired")
	ErrTokenLimitExceeded    = errors.New("token generation limit exceeded")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrForbidden             = errors.New("forbidden")
	ErrSessionExpired        = errors.New("session has expired")
	ErrSubscriptionNotFound  = errors.New("subscription not found")
	ErrTokenNotFound         = errors.New("token not found")
)

// IsUserNotFound проверяет если ошибка - user not found
func IsUserNotFound(err error) bool {
	return errors.Is(err, ErrUserNotFound)
}

// IsInvalidCredentials проверяет если ошибка - invalid credentials
func IsInvalidCredentials(err error) bool {
	return errors.Is(err, ErrInvalidCredentials)
}

// IsUnauthorized проверяет если ошибка - unauthorized
func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}
