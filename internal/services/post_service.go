package service

import (
	"context"
	"errors"
	"strings"

	"weblog/internal/models"
	"weblog/internal/repositories"
)

type PostService struct {
	postRepository *repository.PostRepository
}

func NewPostService(postRepository *repository.PostRepository,) *PostService {
	return &PostService{
		postRepository: postRepository,
	}
}

func (s *PostService) CreatePost(ctx context.Context, title string, content string, image *string, authorID int64, isPrivate bool,
) (*models.Post, error) {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}
	if content == "" {
		return nil, errors.New("content cannot be empty")
	}
	post, err := s.postRepository.CreatePost(ctx, title, content, image, authorID, isPrivate,)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostService) GetPost(ctx context.Context, postID int64,) (*models.Post, error) {
	return s.postRepository.GetPostByID(ctx, postID)
}

func (s *PostService) GetFeed(ctx context.Context, userID int64,) ([]models.Post, error) {
	return s.postRepository.GetPostsForUser(ctx, userID)
}

func (s *PostService) DeletePost(ctx context.Context, postID int64, userID int64,) error {
	return s.postRepository.DeletePost(ctx, postID, userID,)
}