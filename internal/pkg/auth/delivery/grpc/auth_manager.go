package grpc

import (
	"context"
	"errors"

	"github.com/IlyaChgn/voblako/internal/models"
	authinterfaces "github.com/IlyaChgn/voblako/internal/pkg/auth"
	"github.com/IlyaChgn/voblako/internal/pkg/auth/delivery/grpc/protobuf"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthManager struct {
	protobuf.UnimplementedAuthServer

	sessionManager authinterfaces.SessionManager
	authStorage    authinterfaces.AuthRepository
}

func NewAuthManager(
	manager authinterfaces.SessionManager, storage authinterfaces.AuthRepository,
) *AuthManager {
	return &AuthManager{
		sessionManager: manager,
		authStorage:    storage,
	}
}

func (m *AuthManager) CreateUser(ctx context.Context, newUser *protobuf.NewUser) (*protobuf.User, error) {
	user, err := m.authStorage.CreateUser(ctx, newUser.Email, newUser.Password)
	if err != nil {
		if errors.Is(err, models.UserAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "%s", err.Error())
		}

		return nil, err
	}

	return convertUser(user), nil
}

func (m *AuthManager) CreateSession(ctx context.Context, user *protobuf.FullUserData) (*emptypb.Empty, error) {
	return nil, m.sessionManager.CreateSession(ctx, user.SessionID, &models.User{
		ID:    uint(user.User.ID),
		Email: user.User.Email,
	})
}

func (m *AuthManager) Logout(ctx context.Context, session *protobuf.SessionData) (*emptypb.Empty, error) {
	err := m.sessionManager.RemoveSession(ctx, session.SessionID)
	if err != nil {
		if errors.Is(err, models.SessionNotExistsError) {
			return nil, status.Errorf(codes.AlreadyExists, "%s", err.Error())
		}

		return nil, err
	}

	return nil, err
}

func (m *AuthManager) GetUserByEmail(ctx context.Context, email *protobuf.EmailData) (*protobuf.User, error) {
	user, err := m.authStorage.GetUserByEmail(ctx, email.Email)
	return convertUser(user), err
}

func (m *AuthManager) GetCurrentUser(ctx context.Context, session *protobuf.SessionData) (*protobuf.User, error) {
	user, exists := m.sessionManager.GetSession(ctx, session.SessionID)
	if !exists {
		return &protobuf.User{IsAuth: false}, nil
	}

	currUser := convertUser(user)
	currUser.IsAuth = true
	return currUser, nil
}

func convertUser(user *models.User) *protobuf.User {
	if user == nil {
		return nil
	}

	return &protobuf.User{
		ID:           uint32(user.ID),
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}
}
