package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"weblog/internal/models"
	"weblog/internal/repositories"
)

type SessionService struct {
	sessionRepository *repository.SessionRepository
}

func NewSessionService(sessionRepository *repository.SessionRepository,) *SessionService {
	return &SessionService{
		sessionRepository: sessionRepository,
	}
}

func (s *SessionService) CreateSession(ctx context.Context, userID int64,) (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	token := hex.EncodeToString(randomBytes)
	expiresAt := time.Now().Add(24 * time.Hour)
	err = s.sessionRepository.CreateSession(ctx, token, userID, expiresAt,)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *SessionService) GetSession(ctx context.Context, token string,) (*models.Session, error) {
	return s.sessionRepository.GetSessionByToken(ctx, token,)
}

func (s *SessionService) DeleteSession(ctx context.Context, token string,) error {
	return s.sessionRepository.DeleteSession(ctx, token,)
}