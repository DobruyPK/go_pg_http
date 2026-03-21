package user

import "errors"

var (
	ErrInvalidUserName = errors.New("Invalid user name")
	ErrUserNotFound    = errors.New("User not found")
	ErrUserExists      = errors.New("User already exists")
)
