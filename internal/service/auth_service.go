package service

import (
	"context"
	"errors"
	"finance-tracker/internal/models"
	"finance-tracker/internal/repository"
	"finance-tracker/pkg/hash"
	"time"
)

type AuthService struct {
	repo   *repository.UserRepo
	hasher *hash.SHA256Hasher
}

func NewAuthService(repo *repository.UserRepo, hasher *hash.SHA256Hasher) *AuthService {
	return &AuthService{
		repo:   repo,
		hasher: hasher,
	}
}

func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest) error {
	if req.Name == "" || req.Email == "" {
		return errors.New("все поля обязательны")
	}
	if len(req.Password) < 4 {
		return errors.New("пароль слишком короткий")
	}
	if req.Password != req.ConfirmPassword {
		return errors.New("пароли не совпадают")
	}

	if _, err := s.repo.GetByEmail(ctx, req.Email); err == nil {
		// ИСПРАВЛЕНО: Сообщение на русском
		return errors.New("пользователь с таким email уже существует")
	}

	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: s.hasher.Hash(req.Password),
		CreatedAt:    time.Now(),
	}

	return s.repo.Create(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (*models.User, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	if user.PasswordHash != s.hasher.Hash(req.Password) {
		return nil, errors.New("неверный пароль")
	}

	return user, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *AuthService) UpdatePassword(ctx context.Context, userID int, req models.UpdatePasswordRequest) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.PasswordHash != s.hasher.Hash(req.OldPassword) {
		return errors.New("старый пароль неверен")
	}

	if len(req.NewPassword) < 4 {
		return errors.New("новый пароль слишком короткий")
	}

	newHash := s.hasher.Hash(req.NewPassword)
	return s.repo.UpdatePassword(ctx, userID, newHash)
}
