package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"weblog/internal/models"
)

type CommentRepository struct {
	db *pgx.Conn
}

func NewCommentRepository(db *pgx.Conn) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

func (r *CommentRepository) CreateComment(ctx context.Context, postID int64, authorID int64, content string,) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.QueryRow(ctx,
		`INSERT INTO comments (post_id, author_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, post_id, author_id, content, created_at`,
		postID, authorID, content,).Scan(&comment.ID, &comment.PostID, &comment.AuthorID, &comment.Content, &comment.CreatedAt,)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *CommentRepository) GetCommentsByPostID(ctx context.Context, postID int64,) ([]models.Comment, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, post_id, author_id, content, created_at
		FROM comments
		WHERE post_id = $1
		ORDER BY created_at ASC`,
		postID,)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.AuthorID, &comment.Content, &comment.CreatedAt,)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}