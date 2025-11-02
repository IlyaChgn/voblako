package repository

import (
	"context"
	"errors"

	"github.com/IlyaChgn/voblako/internal/models"
	authinterface "github.com/IlyaChgn/voblako/internal/pkg/auth"
	"github.com/IlyaChgn/voblako/internal/pkg/server/dbinit"
	"github.com/IlyaChgn/voblako/internal/pkg/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type authStorage struct {
	pool dbinit.PostgresPool
}

func NewAuthStorage(pool dbinit.PostgresPool) authinterface.AuthRepository {
	return &authStorage{pool: pool}
}

func (s *authStorage) CreateUser(ctx context.Context, email, password string) (*models.User, error) {
	var user models.User

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	line := tx.QueryRow(ctx, CreateUserQuery, email, utils.HashPassword(password))
	if err := line.Scan(&user.ID, &user.Email); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return nil, models.UserAlreadyExists
		}

		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *authStorage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	line := s.pool.QueryRow(ctx, GetUserByEmailQuery, email)
	if err := line.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
