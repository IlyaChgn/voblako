package usecases

import (
	"context"
	"testing"

	"github.com/IlyaChgn/voblako/internal/models"
	"github.com/IlyaChgn/voblako/internal/pkg/auth/delivery/grpc/protobuf"
	"github.com/IlyaChgn/voblako/internal/pkg/auth/usecases/mocks"
	"github.com/IlyaChgn/voblako/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestAuthUsecases_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthClient := mocks.NewMockAuthClient(ctrl)
	au := NewAuthUsecases(mockAuthClient)

	loginData := &models.LoginData{
		Email:    "test@example.com",
		Password: "password",
	}

	user := &protobuf.User{
		ID:           1,
		Email:        loginData.Email,
		PasswordHash: utils.HashPassword(loginData.Password),
	}

	mockAuthClient.EXPECT().GetUserByEmail(gomock.Any(), &protobuf.EmailData{Email: loginData.Email}).Return(user, nil)
	mockAuthClient.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return(&emptypb.Empty{}, nil)

	fullUserData, err := au.Login(context.Background(), loginData)

	assert.NoError(t, err)
	assert.NotNil(t, fullUserData)
	assert.Equal(t, loginData.Email, fullUserData.User.Email)
}

func TestAuthUsecases_Signup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthClient := mocks.NewMockAuthClient(ctrl)
	au := NewAuthUsecases(mockAuthClient)

	signupData := &models.SignupData{
		Email:          "test@example.com",
		Password:       "password",
		PasswordRepeat: "password",
	}

	newUser := &protobuf.User{
		ID:    1,
		Email: signupData.Email,
	}

	mockAuthClient.EXPECT().CreateUser(gomock.Any(), &protobuf.NewUser{Email: signupData.Email, Password: signupData.Password}).Return(newUser, nil)
	mockAuthClient.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return(&emptypb.Empty{}, nil)

	fullUserData, err := au.Signup(context.Background(), signupData)

	assert.NoError(t, err)
	assert.NotNil(t, fullUserData)
	assert.Equal(t, signupData.Email, fullUserData.User.Email)
}

func TestAuthUsecases_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthClient := mocks.NewMockAuthClient(ctrl)
	au := NewAuthUsecases(mockAuthClient)

	sessionID := "session123"

	mockAuthClient.EXPECT().Logout(gomock.Any(), &protobuf.SessionData{SessionID: sessionID}).Return(&emptypb.Empty{}, nil)

	err := au.Logout(context.Background(), sessionID)

	assert.NoError(t, err)
}

func TestAuthUsecases_CheckAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthClient := mocks.NewMockAuthClient(ctrl)
	au := NewAuthUsecases(mockAuthClient)

	sessionID := "session123"
	user := &protobuf.User{ID: 1, Email: "test@example.com", IsAuth: true}

	mockAuthClient.EXPECT().GetCurrentUser(gomock.Any(), &protobuf.SessionData{SessionID: sessionID}).Return(user, nil)

	retrievedUser, isAuth := au.CheckAuth(context.Background(), sessionID)

	assert.True(t, isAuth)
	assert.NotNil(t, retrievedUser)
	assert.Equal(t, user.Email, retrievedUser.Email)
}

func TestAuthUsecases_CheckAuth_NotAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthClient := mocks.NewMockAuthClient(ctrl)
	au := NewAuthUsecases(mockAuthClient)

	sessionID := "session123"
	user := &protobuf.User{IsAuth: false}

	mockAuthClient.EXPECT().GetCurrentUser(gomock.Any(), &protobuf.SessionData{SessionID: sessionID}).Return(user, nil)

	retrievedUser, isAuth := au.CheckAuth(context.Background(), sessionID)

	assert.False(t, isAuth)
	assert.Nil(t, retrievedUser)
}

func TestAuthUsecases_Signup_PasswordMismatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthClient := mocks.NewMockAuthClient(ctrl)
	au := NewAuthUsecases(mockAuthClient)

	signupData := &models.SignupData{
		Email:          "test@example.com",
		Password:       "password",
		PasswordRepeat: "wrongpassword",
	}

	fullUserData, err := au.Signup(context.Background(), signupData)

	assert.Error(t, err)
	assert.Nil(t, fullUserData)
	assert.Equal(t, models.PasswordsNotMatch, err)
}

func TestAuthUsecases_Signup_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthClient := mocks.NewMockAuthClient(ctrl)
	au := NewAuthUsecases(mockAuthClient)

	signupData := &models.SignupData{
		Email:          "test@example.com",
		Password:       "password",
		PasswordRepeat: "password",
	}

	st := status.New(codes.AlreadyExists, "user already exists")
	mockAuthClient.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil, st.Err())

	fullUserData, err := au.Signup(context.Background(), signupData)

	assert.Error(t, err)
	assert.Nil(t, fullUserData)
	assert.Equal(t, models.UserAlreadyExists, err)
}
