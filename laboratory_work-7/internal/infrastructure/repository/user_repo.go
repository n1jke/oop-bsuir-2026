package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/application"
)

type UserRepoSQL struct {
	db *pgxpool.Pool
}

const (
	createUser = `
		INSERT INTO users(id, login, password_hash, rating)
		VALUES ($1, $2, $3, 0)
	`

	getUserByLogin = `
		SELECT id, login, password_hash, rating FROM users
		WHERE login = $1
	`

	getUserByID = `
		SELECT id, login, password_hash, rating FROM users
		WHERE id = $1
	`

	getAllUsers = `
		SELECT id, login, rating FROM users
	`
)

func (r UserRepoSQL) Add(ctx context.Context, login, password string) (uuid.UUID, error) {
	q := GetQuerier(ctx, r.db)

	id := uuid.New()
	_, err := q.Exec(ctx, createUser, id, login, password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return uuid.Nil, application.ErrAlreadyExist
		}

		return uuid.Nil, err
	}

	return id, nil
}

func (r UserRepoSQL) GetByLogin(ctx context.Context, login string) (application.UserRepoDTO, error) {
	q := GetQuerier(ctx, r.db)

	var dto application.UserRepoDTO
	row := q.QueryRow(ctx, getUserByLogin, login)
	if err := row.Scan(&dto.ID, &dto.Name, &dto.PasswordHash, &dto.Rating); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return application.UserRepoDTO{}, application.ErrNotFound
		}

		return application.UserRepoDTO{}, err
	}

	return dto, nil
}

func (r UserRepoSQL) GetByID(ctx context.Context, userID uuid.UUID) (application.UserRepoDTO, error) {
	q := GetQuerier(ctx, r.db)

	var dto application.UserRepoDTO
	row := q.QueryRow(ctx, getUserByID, userID)
	if err := row.Scan(&dto.ID, &dto.Name, &dto.PasswordHash, &dto.Rating); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return application.UserRepoDTO{}, application.ErrNotFound
		}

		return application.UserRepoDTO{}, err
	}

	return dto, nil
}

func (r UserRepoSQL) GetAll(ctx context.Context) ([]*application.UserRepoDTO, error) {
	q := GetQuerier(ctx, r.db)

	rows, err := q.Query(ctx, getAllUsers)
	if err != nil {
		return nil, err
	}

	users := make([]*application.UserRepoDTO, 0)

	for rows.Next() {
		var dto application.UserRepoDTO
		if err := rows.Scan(&dto.ID, &dto.Name, &dto.Rating); err != nil {
			return nil, err
		}

		users = append(users, &dto)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func NewUserRepoSQL(db *pgxpool.Pool) *UserRepoSQL {
	return &UserRepoSQL{db: db}
}
