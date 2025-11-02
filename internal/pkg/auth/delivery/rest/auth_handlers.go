package rest

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/IlyaChgn/voblako/internal/models"
	authinterface "github.com/IlyaChgn/voblako/internal/pkg/auth"
	"github.com/IlyaChgn/voblako/internal/pkg/server/delivery/responses"
)

const sessionDuration = 30 * 24 * time.Hour

type AuthHandler struct {
	usecases authinterface.AuthUsecases
}

func NewAuthHandler(usecases authinterface.AuthUsecases) *AuthHandler {
	return &AuthHandler{usecases: usecases}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, _ := r.Cookie("session_id")
	if session != nil {
		if _, exists := h.usecases.CheckAuth(ctx, session.Value); exists {
			responses.SendErrResponse(w, responses.StatusBadRequest, responses.ErrAuthorized)
			return
		}
	}

	var loginData *models.LoginData
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		responses.SendErrResponse(w, responses.StatusBadRequest, responses.ErrBadJSON)
		return
	}

	user, err := h.usecases.Login(ctx, loginData)
	if err != nil {
		if errors.Is(err, models.PasswordsNotMatch) || errors.Is(err, models.UserNotExists) {
			responses.SendErrResponse(w, responses.StatusBadRequest, responses.ErrWrongCredentials)
			return
		}

		log.Println(err)
		responses.SendErrResponse(w, responses.StatusInternalServerError, responses.ErrInternalServer)
		return
	}

	newSession := createSession(user.SessionID)
	http.SetCookie(w, newSession)
	responses.SendOkResponse(w, &models.AuthData{
		User: models.User{
			ID:    user.ID,
			Email: user.Email,
		},
		IsAuth: true,
	})
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, _ := r.Cookie("session_id")
	if session != nil {
		if _, exists := h.usecases.CheckAuth(ctx, session.Value); exists {
			responses.SendErrResponse(w, responses.StatusBadRequest, responses.ErrAuthorized)

			return
		}
	}

	var signupData *models.SignupData
	err := json.NewDecoder(r.Body).Decode(&signupData)
	if err != nil {
		responses.SendErrResponse(w, responses.StatusBadRequest, responses.ErrBadJSON)
		return
	}

	user, err := h.usecases.Signup(ctx, signupData)
	if err != nil {
		switch {
		case errors.Is(err, models.PasswordsNotMatch):
			responses.SendErrResponse(w, responses.StatusBadRequest, responses.ErrDoNotMatch)
		case errors.Is(err, models.IncorrectPasswordLen):
			responses.SendErrResponse(w, responses.StatusBadRequest, responses.ErrWrongPasswordFormat)
		case errors.Is(err, models.UserAlreadyExists):
			responses.SendErrResponse(w, responses.StatusBadRequest, responses.ErrAlreadyExists)
		default:
			log.Println(err)
			responses.SendErrResponse(w, responses.StatusInternalServerError, responses.ErrInternalServer)
		}
		return
	}

	newSession := createSession(user.SessionID)
	http.SetCookie(w, newSession)
	responses.SendOkResponse(w, &models.AuthData{
		User: models.User{
			ID:    user.ID,
			Email: user.Email,
		},
		IsAuth: true,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, _ := r.Cookie("session_id")
	if session == nil {
		responses.SendErrResponse(w, responses.StatusUnauthorized, responses.ErrNotAuthorized)
		return
	}

	err := h.usecases.Logout(ctx, session.Value)
	if err != nil {
		if errors.Is(err, models.SessionNotExistsError) {
			responses.SendErrResponse(w, responses.StatusUnauthorized, responses.ErrNotAuthorized)
			return
		}

		log.Println(err)
		responses.SendErrResponse(w, responses.StatusInternalServerError, responses.ErrInternalServer)
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
	responses.SendOkResponse(w, nil)
}

func (h *AuthHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, _ := r.Cookie("session_id")
	if session == nil {
		responses.SendOkResponse(w, &models.AuthData{IsAuth: false})
		return
	}

	user, isAuth := h.usecases.CheckAuth(ctx, session.Value)

	if !isAuth {
		responses.SendOkResponse(w, &models.AuthData{IsAuth: false})
		return
	}

	responses.SendOkResponse(w, &models.AuthData{IsAuth: true, User: *user})
}

func createSession(sessionID string) *http.Cookie {
	return &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().Add(sessionDuration),
		HttpOnly: true,
	}
}
