package user

import (
	"git.garena.com/shaoyihong/go-entry-task/common/pb"
	"git.garena.com/shaoyihong/go-entry-task/tcpserver/common"
	"github.com/go-redis/redis/v8"
)

type IUserRepository interface {
	GetByUsername(username string, fromCache bool) (*pb.User, error)
	UpdateNickname(username, nickname string) error
	UpdateProfileImage(username, imageUrl string) error
	Insert(username, password, nickname, imageUrl string) (*pb.User, error)
}

type UserRepository struct {
	dao   IUserDAO
	cache IUserCache
}

func NewUserRepository(database common.Database, redis *redis.Client) IUserRepository {
	return &UserRepository{
		dao:   newUserDAO(database),
		cache: newUserCache(redis),
	}
}
func (repo *UserRepository) GetByUsername(username string, fromCache bool) (*pb.User, error) {
	if fromCache {
		user := repo.cache.getCacheUser(username)
		if user != nil {
			return user, nil
		}
	}
	user, err := repo.dao.getByUsername(username)
	if err != nil {
		return nil, err
	}
	repo.cache.setCacheUser(user)
	return user, nil
}

func (repo *UserRepository) UpdateNickname(username, nickname string) error {
	return repo.dao.updateNickname(username, nickname)
}

func (repo *UserRepository) UpdateProfileImage(username, imageUrl string) error {
	return repo.dao.updateProfileImage(username, imageUrl)
}

func (repo *UserRepository) Insert(username, password, nickname, imageUrl string) (*pb.User, error) {
	user, err := repo.dao.insert(username, password, nickname, imageUrl)
	if err == nil && user != nil {
		repo.cache.setCacheUser(user)
	}
	return user, err
}
