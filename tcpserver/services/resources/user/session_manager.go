package user

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"time"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"github.com/go-redis/redis/v8"
)

const (
	tokenExpiryTimeInHours = 12
	tokenSuffix            = "-token"
)

type ISessionManager interface {
	SetCacheToken(username string) (string, error)
	GetCacheToken(username string) (string, error)
	GetCacheUsername(token string) (string, error)
	DeleteCacheToken(username string)
}

type SessionManager struct {
	redis *redis.Client
}

func newSessionManager(redis *redis.Client) ISessionManager {
	return &SessionManager{redis: redis}
}

func createSessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}

	return base64.URLEncoding.EncodeToString(b)
}

func createKeyFromUsername(username string) string {
	return username + tokenSuffix
}

func generateUsernameFromKey(key string) string {
	return key[:len(key)-len(tokenSuffix)]
}

func (manager *SessionManager) SetCacheToken(username string) (string, error) {
	token := createSessionID()
	key := createKeyFromUsername(username)

	if err := manager.redis.Set(context.Background(), key, token, time.Hour*tokenExpiryTimeInHours).Err(); err != nil {
		logger.ErrorLogger.Println("Failed to set cache token:", err)

		return "", err
	}

	if err := manager.redis.Set(context.Background(), token, key, time.Hour*tokenExpiryTimeInHours).Err(); err != nil {
		logger.ErrorLogger.Println("Failed to set reverse cache token:", err)

		return "", err
	}

	return token, nil
}

func (manager *SessionManager) GetCacheToken(username string) (string, error) {
	token, err := manager.redis.Get(context.Background(), createKeyFromUsername(username)).Result()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (manager *SessionManager) GetCacheUsername(token string) (string, error) {
	key, err := manager.redis.Get(context.Background(), token).Result()
	if err != nil {
		return "", err
	}

	return generateUsernameFromKey(key), nil
}

func (manager *SessionManager) DeleteCacheToken(username string) {
	manager.redis.Del(context.Background(), createKeyFromUsername(username))
}
