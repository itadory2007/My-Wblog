package service

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"weblog/internal/models"
	"weblog/internal/repositories"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository,) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) Register(ctx context.Context, username string, password string,) (*models.User, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	if password == "" {
		return nil, errors.New("password cannot be empty")
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost,)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.CreateUser(ctx, username, string(passwordHash),)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Login(ctx context.Context, username string, password string,) (*models.User, error) {
	username = strings.TrimSpace(username)
	user, err := s.userRepository.GetUserByUsername(ctx, username,)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password),)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}
	return user, nil
}