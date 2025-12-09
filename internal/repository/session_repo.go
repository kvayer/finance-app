package repository

import (
	"context"
	"database/sql"
	"errors"
	"finance-tracker/internal/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const sessionsTable = "sessions"

type SessionRepo struct {
	db *sqlx.DB
}

func NewSessionRepo(db *sqlx.DB) *SessionRepo {
	return &SessionRepo{db: db}
}

func (r *SessionRepo) Create(ctx context.Context, s *models.Session) error {
	query := fmt.Sprintf("INSERT INTO %s (token, user_id, expiry) VALUES ($1, $2, $3)", sessionsTable)
	_, err := r.db.ExecContext(ctx, query, s.Token, s.UserID, s.Expiry)
	return err
}

func (r *SessionRepo) GetByToken(ctx context.Context, token string) (*models.Session, error) {
	var s models.Session
	query := fmt.Sprintf("SELECT * FROM %s WHERE token = $1", sessionsTable)
	if err := r.db.GetContext(ctx, &s, query, token); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Session not found
		}
		return nil, err
	}
	return &s, nil
}

func (r *SessionRepo) Delete(ctx context.Context, token string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE token = $1", sessionsTable)
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}
