package user

import (
	"git.garena.com/shaoyihong/go-entry-task/common/pb"
	"git.garena.com/shaoyihong/go-entry-task/tcpserver/common"
	"github.com/go-redis/redis/v8"
)

type IUserRepository interface {
	GetByUsername(username string, fromCache bool) (*pb.User, error)
	UpdateNickname(username, nickname string) error
	UpdateProfileImage(username, imageURL string) error
	Insert(username, password, nickname, imageURL string) (*pb.User, error)
}

type Repository struct {
	dao   IUserDAO
	cache IUserCache
}

func NewUserRepository(database common.Database, redis *redis.Client) IUserRepository {
	return &Repository{
		dao:   newUserDAO(database),
		cache: newUserCache(redis),
	}
}

func (repo *Repository) GetByUsername(username string, fromCache bool) (*pb.User, error) {
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

func (repo *Repository) UpdateNickname(username, nickname string) error {
	return repo.dao.updateNickname(username, nickname)
}

func (repo *Repository) UpdateProfileImage(username, imageURL string) error {
	return repo.dao.updateProfileImage(username, imageURL)
}

func (repo *Repository) Insert(username, password, nickname, imageURL string) (*pb.User, error) {
	user, err := repo.dao.insert(username, password, nickname, imageURL)
	if err == nil && user != nil {
		repo.cache.setCacheUser(user)
	}

	return user, err
}
