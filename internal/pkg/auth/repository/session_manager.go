package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IlyaChgn/voblako/internal/models"
	authinterface "github.com/IlyaChgn/voblako/internal/pkg/auth"

	"github.com/redis/go-redis/v9"
)

const sessionDuration = 30 * 24 * time.Hour

type sessionManager struct {
	client *redis.Client
}

func NewSessionManager(client *redis.Client) authinterface.SessionManager {
	return &sessionManager{
		client: client,
	}
}

func (manager *sessionManager) CreateSession(ctx context.Context, sessionID string, user *models.User) error {
	rawUser, err := json.Marshal(user)
	if err != nil {
		return models.MarshallingSessionError
	}

	err = manager.client.Set(ctx, sessionID, rawUser, sessionDuration).Err()
	if err != nil {
		return models.AddToRedisError
	}

	return nil
}

func (manager *sessionManager) RemoveSession(ctx context.Context, sessionID string) error {
	if _, exists := manager.GetSession(ctx, sessionID); !exists {
		return models.SessionNotExistsError
	}

	_, err := manager.client.Del(ctx, sessionID).Result()
	if err != nil {
		return models.DeleteFromRedisError
	}

	return nil
}

func (manager *sessionManager) GetSession(ctx context.Context, sessionID string) (*models.User, bool) {
	rawUser, _ := manager.client.Get(ctx, sessionID).Result()

	var user *models.User
	if err := json.Unmarshal([]byte(rawUser), &user); err != nil {
		return nil, false
	}

	return user, user != nil
}
