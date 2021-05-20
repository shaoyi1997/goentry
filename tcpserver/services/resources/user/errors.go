package user

import "errors"

var (
	errUsernameAlreadyExists = errors.New("username already exists")
	errUsernameNotFound      = errors.New("username not found")
	errEmptyUsername         = errors.New("username cannot be empty")
	errEmptyPassword         = errors.New("password cannot be empty")
	errTooShortPassword      = errors.New("password is too short")
)
