package session

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"finance-tracker/internal/models"
	"finance-tracker/internal/repository"
	"time"
)

type Manager struct {
	repo *repository.SessionRepo
}

func NewManager(repo *repository.SessionRepo) *Manager {
	return &Manager{repo: repo}
}

func (m *Manager) Create(ctx context.Context, userID int) (string, error) {
	token := generateToken()
	expiry := time.Now().Add(24 * time.Hour) // Сессия на сутки

	s := &models.Session{
		Token:  token,
		UserID: userID,
		Expiry: expiry,
	}

	if err := m.repo.Create(ctx, s); err != nil {
		return "", err
	}

	return token, nil
}

func (m *Manager) Check(ctx context.Context, token string) (int, error) {
	s, err := m.repo.GetByToken(ctx, token)
	if err != nil {
		return 0, err
	}
	if s == nil {
		return 0, nil // Invalid token
	}

	if time.Now().After(s.Expiry) {
		_ = m.repo.Delete(ctx, token)
		return 0, nil // Expired
	}

	return s.UserID, nil
}

func generateToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
