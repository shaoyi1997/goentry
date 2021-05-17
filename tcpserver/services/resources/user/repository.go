package user

import (
	"git.garena.com/shaoyihong/go-entry-task/common/pb"
	"git.garena.com/shaoyihong/go-entry-task/tcpserver/common"
)

type IUserRepository interface {
	GetByUsername(username string) (*pb.User, error)
	UpdateNickname(username, nickname string) error
	UpdateProfileImage(username, imageUrl string) error
	Insert(username, password, nickname, imageUrl string) (*pb.User, error)
}

type UserRepository struct {
	dao IUserDAO
}

func NewUserRepository(database common.Database) IUserRepository {
	return &UserRepository{
		dao: newUserDAO(database),
	}
}
func (repo *UserRepository) GetByUsername(username string) (*pb.User, error) {
	return repo.dao.getByUsername(username)
}

func (repo *UserRepository) UpdateNickname(username, nickname string) error {
	return repo.dao.updateNickname(username, nickname)
}

func (repo *UserRepository) UpdateProfileImage(username, imageUrl string) error {
	return repo.dao.updateProfileImage(username, imageUrl)
}

func (repo *UserRepository) Insert(username, password, nickname, imageUrl string) (*pb.User, error) {
	return repo.dao.insert(username, password, nickname, imageUrl)
}
