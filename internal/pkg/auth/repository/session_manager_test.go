package repository

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/IlyaChgn/voblako/internal/models"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestSessionManager_CreateSession(t *testing.T) {
	client, mock := redismock.NewClientMock()
	sm := NewSessionManager(client)

	user := &models.User{ID: 1, Email: "test@example.com"}
	sessionID := "session123"

	rawUser, _ := json.Marshal(user)
	mock.ExpectSet(sessionID, rawUser, sessionDuration).SetVal("")

	err := sm.CreateSession(context.Background(), sessionID, user)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSessionManager_RemoveSession(t *testing.T) {
	client, mock := redismock.NewClientMock()
	sm := NewSessionManager(client)

	user := &models.User{ID: 1, Email: "test@example.com"}
	sessionID := "session123"

	rawUser, _ := json.Marshal(user)
	mock.ExpectGet(sessionID).SetVal(string(rawUser))
	mock.ExpectDel(sessionID).SetVal(1)

	err := sm.RemoveSession(context.Background(), sessionID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSessionManager_GetSession(t *testing.T) {
	client, mock := redismock.NewClientMock()
	sm := NewSessionManager(client)

	user := &models.User{ID: 1, Email: "test@example.com"}
	sessionID := "session123"

	rawUser, _ := json.Marshal(user)
	mock.ExpectGet(sessionID).SetVal(string(rawUser))

	retrievedUser, exists := sm.GetSession(context.Background(), sessionID)

	assert.True(t, exists)
	assert.Equal(t, user, retrievedUser)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSessionManager_GetSession_NotFound(t *testing.T) {
	client, mock := redismock.NewClientMock()
	sm := NewSessionManager(client)

	sessionID := "session123"

	mock.ExpectGet(sessionID).RedisNil()

	retrievedUser, exists := sm.GetSession(context.Background(), sessionID)

	assert.False(t, exists)
	assert.Nil(t, retrievedUser)
	assert.NoError(t, mock.ExpectationsWereMet())
}
