package user

import (
	"errors"

	"git.garena.com/shaoyihong/go-entry-task/common/pb"
	"git.garena.com/shaoyihong/go-entry-task/tcpserver/common"
)

type UserService struct {
	repo      IUserRepository
	hasher    IPasswordHasher
	validator IUserValidator
}

func NewUserService(database common.Database) *UserService {
	return &UserService{
		repo:      NewUserRepository(database),
		hasher:    newPasswordHasher(),
		validator: newUserValidator(),
	}
}
func (service *UserService) GetByUsername(username string) (*pb.User, error) {
	return service.repo.GetByUsername(username)
}

func (service *UserService) UpdateNickname(username, nickname string) error {
	return service.repo.UpdateNickname(username, nickname)
}

func (service *UserService) UpdateProfileImage(username, imageUrl string) error {
	return service.repo.UpdateProfileImage(username, imageUrl)
}

func (service *UserService) Register(username, password, nickname, imageUrl string) (*pb.User, error) {
	err := service.validator.ValidateLogin(username, password)
	if err != nil {
		return nil, err
	}
	hashedPassword, err := service.hasher.hash(password)
	if err != nil {
		return nil, err
	}
	return service.repo.Insert(username, hashedPassword, nickname, imageUrl)
}

func (service *UserService) Login(username, password string) (*pb.User, error) {
	err := service.validator.ValidateLogin(username, password)
	if err != nil {
		return nil, err
	}
	user, err := service.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("username not found")
	}
	isValidPassword := service.hasher.comparePasswords(*user.Password, password)
	if isValidPassword {
		return user, nil
	} else {
		return nil, errors.New("password is invalid")
	}
}
