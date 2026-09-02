package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"weblog/internal/models"
)

type PostAccessRepository struct {
	db *pgx.Conn
}

func NewPostAccessRepository(db *pgx.Conn) *PostAccessRepository {
	return &PostAccessRepository{
		db: db,
	}
}

func (r *PostAccessRepository) GrantAccess(ctx context.Context, postID int64, userID int64,) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO post_access (post_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (post_id, user_id) DO NOTHING`,
		postID, userID,)
	return err
}

func (r *PostAccessRepository) RevokeAccess(ctx context.Context, postID int64, userID int64,) error {
	commandTag, err := r.db.Exec(ctx,
		`DELETE FROM post_access
		WHERE post_id = $1 AND user_id = $2`,
		postID, userID,)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *PostAccessRepository) GetUsersWithAccess(ctx context.Context, postID int64,) ([]models.User, error) {
	rows, err := r.db.Query(ctx,
		`SELECT u.id, u.username, u.password_hash, u.created_at
		FROM users u
		INNER JOIN post_access pa
			ON u.id = pa.user_id
		WHERE pa.post_id = $1
		ORDER BY u.username`,
		postID,)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt,)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}