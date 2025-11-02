package auth

import (
	"context"

	"github.com/IlyaChgn/voblako/internal/models"
)

type SessionManager interface {
	CreateSession(ctx context.Context, sessionID string, user *models.User) error
	RemoveSession(ctx context.Context, sessionID string) error
	GetSession(ctx context.Context, sessionID string) (*models.User, bool)
}

type AuthRepository interface {
	CreateUser(ctx context.Context, email, password string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type AuthUsecases interface {
	Login(ctx context.Context, data *models.LoginData) (*models.FullUserData, error)
	Signup(ctx context.Context, data *models.SignupData) (*models.FullUserData, error)
	Logout(ctx context.Context, sessionID string) error
	CheckAuth(ctx context.Context, sessionID string) (*models.User, bool)
}
