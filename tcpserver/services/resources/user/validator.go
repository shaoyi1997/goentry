package user

const (
	minimumPasswordLength = 4
)

type IUserValidator interface {
	ValidateNonEmptyUsernamePassword(username, password string) error
	ValidateLogout(username string) error
	ValidateRegister(username, password string) error
}

type userValidator struct{}

func newUserValidator() IUserValidator {
	return &userValidator{}
}

func (validator *userValidator) ValidateNonEmptyUsernamePassword(username, password string) error {
	var err error
	if username == "" {
		err = errEmptyUsername
	} else if password == "" {
		err = errEmptyPassword
	}

	return err
}

func (validator *userValidator) ValidateRegister(username, password string) error {
	err := validator.ValidateNonEmptyUsernamePassword(username, password)
	if err != nil {
		return err
	}

	if len(password) < minimumPasswordLength {
		return errTooShortPassword
	}

	return nil
}

func (validator *userValidator) ValidateLogout(username string) error {
	if username == "" {
		return errEmptyUsername
	}

	return nil
}
