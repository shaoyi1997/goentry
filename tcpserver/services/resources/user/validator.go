package user

import "errors"

type IUserValidator interface {
	ValidateLogin(username, password string) error
	ValidateRegistration(username, password string) error
}

type userValidator struct {
}

func newUserValidator() IUserValidator {
	return &userValidator{}
}

func validateNonEmptyUserNamePassword(username, password string) error {
	var err error
	if username == "" {
		err = errors.New("username cannot be empty")
	} else if password == "" {
		err = errors.New("password cannot be empty")
	}
	return err
}

func (validator *userValidator) ValidateLogin(username, password string) error {
	err := validateNonEmptyUserNamePassword(username, password)
	if err != nil {
		return err
	}
	return nil
}

func (validator *userValidator) ValidateRegistration(username, password string) error {
	err := validateNonEmptyUserNamePassword(username, password)
	if err != nil {
		return err
	}
	return nil
}
