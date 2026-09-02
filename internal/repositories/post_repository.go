package repository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"weblog/internal/models"
)

type PostRepository struct {
	db *pgx.Conn
}

func NewPostRepository(db *pgx.Conn) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (r *PostRepository) CreatePost(
	ctx context.Context,
	title string,
	content string,
	image *string,
	authorID int64,
	isPrivate bool,
) (*models.Post, error) {
	var post models.Post

	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO posts (
			title,
			content,
			image,
			author_id,
			is_private
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, title, content, image, author_id, is_private, created_at
		`,
		title,
		content,
		image,
		authorID,
		isPrivate,
	).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.Image,
		&post.AuthorID,
		&post.IsPrivate,
		&post.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) GetPostByID(
	ctx context.Context,
	id int64,
) (*models.Post, error) {
	var post models.Post

	err := r.db.QueryRow(
		ctx,
		`
		SELECT id, title, content, image, author_id, is_private, created_at
		FROM posts
		WHERE id = $1
		`,
		id,
	).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.Image,
		&post.AuthorID,
		&post.IsPrivate,
		&post.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) GetPublicPosts(
	ctx context.Context,
) ([]models.Post, error) {
	rows, err := r.db.Query(
		ctx,
		`
		SELECT id, title, content, image, author_id, is_private, created_at
		FROM posts
		WHERE is_private = FALSE
		ORDER BY created_at DESC
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.Image,
			&post.AuthorID,
			&post.IsPrivate,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) GetPostsForUser(
	ctx context.Context,
	userID int64,
) ([]models.Post, error) {
	rows, err := r.db.Query(
		ctx,
		`
		SELECT DISTINCT
			p.id,
			p.title,
			p.content,
			p.image,
			p.author_id,
			p.is_private,
			p.created_at
		FROM posts p
		LEFT JOIN post_access pa
			ON p.id = pa.post_id
		WHERE
			p.is_private = FALSE
			OR p.author_id = $1
			OR pa.user_id = $1
		ORDER BY p.created_at DESC
		`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.Image,
			&post.AuthorID,
			&post.IsPrivate,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) DeletePost(
	ctx context.Context,
	postID int64,
	userID int64,
) error {
	commandTag, err := r.db.Exec(
		ctx,
		`
		DELETE FROM posts
		WHERE id = $1 AND author_id = $2
		`,
		postID,
		userID,
	)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}