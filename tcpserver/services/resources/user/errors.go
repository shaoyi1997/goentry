package user

import "errors"

var (
	usernameAlreadyExistsError = errors.New("username already exists")
	usernameNotFoundError      = errors.New("username not found")
	emptyUsernameError         = errors.New("username cannot be empty")
	emptyPasswordError         = errors.New("password cannot be empty")
	emptyTokenError            = errors.New("token cannot be empty")
	tooShortPasswordError      = errors.New("password is too short")
)
