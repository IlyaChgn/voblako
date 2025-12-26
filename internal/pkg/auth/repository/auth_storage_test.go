package repository

import (
	"context"
	"testing"

	"github.com/IlyaChgn/voblako/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func TestAuthStorage_CreateUser_Success(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	s := NewAuthStorage(mock)

	email := "test@example.com"
	password := "password"

	rows := pgxmock.NewRows([]string{"id", "email"}).AddRow(uint(1), email)

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO public.user").WithArgs(email, pgxmock.AnyArg()).WillReturnRows(rows)
	mock.ExpectCommit()

	user, err := s.CreateUser(context.Background(), email, password)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	if user != nil {
		assert.Equal(t, email, user.Email)
		assert.Equal(t, uint(1), user.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAuthStorage_CreateUser_UserAlreadyExists(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	s := NewAuthStorage(mock)

	email := "test@example.com"
	password := "password"

	pgErr := &pgconn.PgError{
		Code:    "23505", // unique_violation
		Message: "duplicate key value violates unique constraint",
	}

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO public.user").WithArgs(email, pgxmock.AnyArg()).WillReturnError(pgErr)
	mock.ExpectRollback()

	user, err := s.CreateUser(context.Background(), email, password)

	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Equal(t, models.UserAlreadyExists, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAuthStorage_GetUserByEmail(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	s := NewAuthStorage(mock)
	email := "test@example.com"
	rows := pgxmock.NewRows([]string{"id", "email", "password_hash"}).
		AddRow(uint(1), email, "hashed_password")

	mock.ExpectQuery("SELECT u.id, u.email, u.password_hash").WithArgs(email).WillReturnRows(rows)

	user, err := s.GetUserByEmail(context.Background(), email)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	if user != nil {
		assert.Equal(t, email, user.Email)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAuthStorage_GetUserByEmail_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	s := NewAuthStorage(mock)
	email := "test@example.com"

	mock.ExpectQuery("SELECT u.id, u.email, u.password_hash").WithArgs(email).WillReturnError(pgx.ErrNoRows)

	user, err := s.GetUserByEmail(context.Background(), email)

	assert.NoError(t, err)
	assert.Nil(t, user)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
