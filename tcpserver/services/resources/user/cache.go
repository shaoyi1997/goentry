package user

import (
	"context"
	"errors"
	"time"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"git.garena.com/shaoyihong/go-entry-task/common/pb"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/proto"
)

const (
	cacheExpiryTime = 12
)

type IUserCache interface {
	setCacheUser(user *pb.User)
	getCacheUser(username string) *pb.User
}

type Cache struct {
	redis *redis.Client
}

func newUserCache(redis *redis.Client) IUserCache {
	return &Cache{redis: redis}
}

func (cache *Cache) setCacheUser(user *pb.User) {
	encodedUser, err := proto.Marshal(user)
	if err != nil {
		logger.ErrorLogger.Println("Failed to encode user:", err)
	}

	if err = cache.redis.Set(context.Background(), *user.Username, string(encodedUser),
		time.Hour*cacheExpiryTime).Err(); err != nil {
		logger.ErrorLogger.Println("Failed to set user in cache:", err)
	}
}

func (cache *Cache) getCacheUser(username string) *pb.User {
	user := new(pb.User)

	encodedUser, err := cache.redis.Get(context.Background(), username).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
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
