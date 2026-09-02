package service

import (
	"context"
	"errors"
	"strings"

	"weblog/internal/models"
	"weblog/internal/repositories"
)

type PostAccessService struct {
	postAccessRepository *repository.PostAccessRepository
	userRepository       *repository.UserRepository
	postRepository       *repository.PostRepository
}

func NewPostAccessService(postAccessRepository *repository.PostAccessRepository, userRepository *repository.UserRepository,postRepository *repository.PostRepository,) *PostAccessService {
	return &PostAccessService{
		postAccessRepository: postAccessRepository,
		userRepository:       userRepository,
		postRepository:       postRepository,
	}
}

func (s *PostAccessService) GrantAccess(ctx context.Context, postID int64, ownerID int64, username string,) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return errors.New("username cannot be empty")
	}

	post, err := s.postRepository.GetPostByID(ctx, postID)
	if err != nil {
		return err
	}
	if post.AuthorID != ownerID {
		return errors.New("only the post owner can grant access")
	}

	user, err := s.userRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return errors.New("user not found")
	}
	if user.ID == ownerID {
		return errors.New("cannot grant access to yourself")
	}
	return s.postAccessRepository.GrantAccess(ctx, postID, user.ID,)
}

func (s *PostAccessService) RevokeAccess(ctx context.Context, postID int64, ownerID int64, userID int64,) error {
	post, err := s.postRepository.GetPostByID(ctx, postID)
	if err != nil {
		return err
	}
	if post.AuthorID != ownerID {
		return errors.New("only the post owner can revoke access")
	}
	return s.postAccessRepository.RevokeAccess(ctx, postID, userID,)
}

func (s *PostAccessService) GetUsersWithAccess(ctx context.Context, postID int64, ownerID int64,) ([]models.User, error) {
	post, err := s.postRepository.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if post.AuthorID != ownerID {
		return nil, errors.New("only the post owner can view access list")
	}
	return s.postAccessRepository.GetUsersWithAccess(ctx, postID,)
}