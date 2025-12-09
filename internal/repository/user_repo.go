package repository

import (
	"context"
	"database/sql"
	"errors"
	"finance-tracker/internal/models"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const usersTable = "users"

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, u *models.User) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (name, email, password_hash, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id`, usersTable)

	err := r.db.QueryRowContext(ctx, query, u.Name, u.Email, u.PasswordHash, time.Now()).Scan(&u.ID)
	if err != nil {
		return fmt.Errorf("repo: create user: %w", err)
	}
	return nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE email = $1", usersTable)

	if err := r.db.GetContext(ctx, &u, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	var u models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", usersTable)
	if err := r.db.GetContext(ctx, &u, query, id); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) UpdatePassword(ctx context.Context, userID int, newHash string) error {
	query := fmt.Sprintf("UPDATE %s SET password_hash = $1 WHERE id = $2", usersTable)
	_, err := r.db.ExecContext(ctx, query, newHash, userID)
	return err
}
