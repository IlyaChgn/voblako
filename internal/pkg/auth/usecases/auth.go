package usecases

import (
	"context"

	"github.com/IlyaChgn/voblako/internal/models"
	authinterface "github.com/IlyaChgn/voblako/internal/pkg/auth"
	"github.com/IlyaChgn/voblako/internal/pkg/auth/delivery/grpc/protobuf"
	"github.com/IlyaChgn/voblako/internal/pkg/utils"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authUsecases struct {
	client protobuf.AuthClient
}

func NewAuthUsecases(client protobuf.AuthClient) authinterface.AuthUsecases {
	return &authUsecases{client: client}
}

func (uc *authUsecases) Login(ctx context.Context, data *models.LoginData) (*models.FullUserData, error) {
	user, err := uc.client.GetUserByEmail(ctx, &protobuf.EmailData{Email: data.Email})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, models.UserNotExists
	}

	if !utils.CheckPasswordHash(data.Password, user.PasswordHash) {
		return nil, models.PasswordsNotMatch
	}

	sessionID := uuid.NewString()
	_, err = uc.client.CreateSession(ctx, &protobuf.FullUserData{
		User:      user,
		SessionID: sessionID,
	})
	if err != nil {
		return nil, err
	}

	return &models.FullUserData{
		User: models.User{
			ID:    uint(user.ID),
			Email: user.Email,
		},
		SessionID: sessionID,
	}, nil
}

func (uc *authUsecases) Signup(ctx context.Context, data *models.SignupData) (*models.FullUserData, error) {
	if data.Password != data.PasswordRepeat {
		return nil, models.PasswordsNotMatch
	}
	if len(data.Password) < 8 || len(data.Password) > 32 {
		return nil, models.IncorrectPasswordLen
	}

	newUser, err := uc.client.CreateUser(ctx, &protobuf.NewUser{
		Email:    data.Email,
		Password: data.Password,
	})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.AlreadyExists {
			return nil, models.UserAlreadyExists
		}

		return nil, err
	}

	sessionID := uuid.NewString()
	_, err = uc.client.CreateSession(ctx, &protobuf.FullUserData{
		User:      newUser,
		SessionID: sessionID,
	})
	if err != nil {
		return nil, err
	}

	return &models.FullUserData{
		User: models.User{
			ID:    uint(newUser.ID),
			Email: newUser.Email,
		},
		SessionID: sessionID,
	}, nil
}

func (uc *authUsecases) Logout(ctx context.Context, sessionID string) error {
	_, err := uc.client.Logout(ctx, &protobuf.SessionData{SessionID: sessionID})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.AlreadyExists {
			return models.SessionNotExistsError
		}

		return err
	}

	return nil
}

func (uc *authUsecases) CheckAuth(ctx context.Context, sessionID string) (*models.User, bool) {
	user, err := uc.client.GetCurrentUser(ctx, &protobuf.SessionData{SessionID: sessionID})
	if err != nil || !user.IsAuth {
		return nil, false
	}

	return &models.User{
		ID:    uint(user.ID),
		Email: user.Email,
	}, true
}
