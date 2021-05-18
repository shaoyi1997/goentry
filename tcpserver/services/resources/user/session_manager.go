package user

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"time"

	"github.com/go-redis/redis/v8"
)

type ISessionManager interface {
	SetCacheToken(username string) (string, error)
	GetCacheToken(username string) (string, error)
	DeleteCacheToken(username string)
}

type SessionManager struct {
	redis *redis.Client
}

func newSessionManager(redis *redis.Client) ISessionManager {
	return &SessionManager{redis: redis}
}

func createSessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func createKeyFromUsername(username string) string {
	return username + "-token"
}

func (manager *SessionManager) SetCacheToken(username string) (string, error) {
	token := createSessionId()
	if err := manager.redis.Set(context.Background(), createKeyFromUsername(username), token, time.Minute*5).Err(); err != nil {
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

func (manager *SessionManager) DeleteCacheToken(username string) {
	manager.redis.Del(context.Background(), createKeyFromUsername(username))
}
