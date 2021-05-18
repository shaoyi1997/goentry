package user

type IUserValidator interface {
	ValidateNonEmptyUsernamePassword(username, password string) error
	ValidateLogout(username string) error
	ValidateRegister(username, password string) error
}

type userValidator struct {
}

func newUserValidator() IUserValidator {
	return &userValidator{}
}

func (validator *userValidator) ValidateNonEmptyUsernamePassword(username, password string) error {
	var err error
	if username == "" {
		err = emptyUsernameError
	} else if password == "" {
		err = emptyPasswordError
	}
	return err
}

func (validator *userValidator) ValidateRegister(username, password string) error {
	err := validator.ValidateNonEmptyUsernamePassword(username, password)
	if err != nil {
		return err
	}
	if len(password) < 4 {
		return tooShortPasswordError
	}
	return nil
}

func (validator *userValidator) ValidateLogout(username string) error {
	if username == "" {
		return emptyUsernameError
	}
	return nil
}
