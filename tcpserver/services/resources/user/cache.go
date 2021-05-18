package user

import (
	"context"
	"time"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"git.garena.com/shaoyihong/go-entry-task/common/pb"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/proto"
)

type IUserCache interface {
	setCacheToken(token, username string) error
	getCacheToken(token string) (string, error)
	setCacheUser(user *pb.User) error
	getCacheUser(username string) *pb.User
}

type UserCache struct {
	redis *redis.Client
}

func newUserCache(redis *redis.Client) IUserCache {
	return &UserCache{redis: redis}
}

func (cache *UserCache) setCacheToken(token, username string) error {
	if err := cache.redis.Set(context.Background(), token, username, time.Minute*5).Err(); err != nil {
		return err
	}
	return nil
}

func (cache *UserCache) getCacheToken(token string) (string, error) {
	username, err := cache.redis.Get(context.Background(), token).Result()
	if err != nil {
		return "", err
	}
	return username, nil
}

func (cache *UserCache) setCacheUser(user *pb.User) error {
	encodedUser, err := proto.Marshal(user)
	if err != nil {
		logger.ErrorLogger.Println("Failed to encode user:", err)
		return err
	}
	if err = cache.redis.Set(context.Background(), *user.Username, string(encodedUser), time.Hour*12).Err(); err != nil {
		return err
	}

	return nil
}

func (cache *UserCache) getCacheUser(username string) *pb.User {
	user := &pb.User{}
	encodedUser, err := cache.redis.Get(context.Background(), username).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		logger.ErrorLogger.Println("Failed to retrieve user:", err)
		return nil
	}

	err = proto.Unmarshal([]byte(encodedUser), user)
	if err != nil {
		logger.ErrorLogger.Println("Failed to unmarshal user from redis:", err)
		return nil
	}
	return user
}
