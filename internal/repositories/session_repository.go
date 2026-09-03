package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"weblog/internal/models"
)

type SessionRepository struct {
	db *pgx.Conn
}

func NewSessionRepository(db *pgx.Conn) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (r *SessionRepository) CreateSession(ctx context.Context, token string, userID int64, eexpiresAt time.Time,) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO sessions (token, user_id, expires_at)
		VALUES ($1, $2, $3)`,
		token, userID, expiresAt,)
	return err
}

func (r *SessionRepository) GetSessionByToken(ctx context.Context, token string,) (*models.Session, error) {
	var session models.Session
	err := r.db.QueryRow(ctx,
		`SELECT id, token, user_id, expires_at, created_at
		FROM sessions
		WHERE token = $1
		  AND expires_at > NOW()`,
		token,).Scan(&session.ID, &session.Token, &session.UserID, &session.ExpiresAt, &session.CreatedAt,)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) DeleteSession(ctx context.Context, token string,) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM sessions
		WHERE token = $1`,
		token,)
	return err
}