package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"weblog/internal/models"
)

type UserRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, username string, passwordHash string,) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(ctx,
		`INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id, username, password_hash, created_at`,
		username, passwordHash,).Scan(&user.ID, &user.Username, &user.PasswordHash,&user.CreatedAt,)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string,) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(ctx,
		`SELECT id, username, password_hash, created_at
		FROM users
		WHERE username = $1`,
		username,).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt,)
	if err != nil {
		return nil, err
	}
	return &user, nil
}