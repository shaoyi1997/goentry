package common

import (
	"context"
	"fmt"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"git.garena.com/shaoyihong/go-entry-task/tcpserver/config"
	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func InitRedis() *redis.Client {
	redisConfig := config.GetRedisConfig()
	redisClient = redis.NewClient(&redis.Options{ //nolint:exhaustivestruct
		Addr:     fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       0,
		PoolSize: redisConfig.PoolSize,
	})

	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		logger.ErrorLogger.Fatalln("Failed to initialize redis:", err)
	}

	logger.InfoLogger.Println("Redis connection initialised successfully")

	return redisClient
}

func TearDownRedis() {
	logger.InfoLogger.Println("Closing redis connection")

	if err := redisClient.Close(); err != nil {
		logger.ErrorLogger.Println("Failed to close redis:", err)
	}
}
