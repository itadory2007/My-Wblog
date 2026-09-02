package service

import (
	"context"
	"errors"
	"strings"

	"weblog/internal/models"
	"weblog/internal/repositories"
)

type CommentService struct {
	commentRepository *repository.CommentRepository
}

func NewCommentService(commentRepository *repository.CommentRepository,) *CommentService {
	return &CommentService{
		commentRepository: commentRepository,
	}
}

func (s *CommentService) CreateComment(ctx context.Context, postID int64, authorID int64, content string,) (*models.Comment, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, errors.New("comment cannot be empty")
	}

	comment, err := s.commentRepository.CreateComment(ctx, postID, authorID, content,)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *CommentService) GetComments(ctx context.Context, postID int64,) ([]models.Comment, error) {
	return s.commentRepository.GetCommentsByPostID(ctx, postID,)
}