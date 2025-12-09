package repository

import (
	"context"
	"finance-tracker/internal/models"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const expensesTable = "expenses"

type ExpenseRepo struct {
	db *sqlx.DB
}

func NewExpenseRepo(db *sqlx.DB) *ExpenseRepo {
	return &ExpenseRepo{db: db}
}

func (r *ExpenseRepo) Create(ctx context.Context, e *models.Expense) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (user_id, category, amount, description, created_at) 
		VALUES ($1, $2, $3, $4, $5)`, expensesTable)

	// Se o tempo não for definido, usamos o atual
	createdAt := e.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	_, err := r.db.ExecContext(ctx, query, e.UserID, e.Category, e.Amount, e.Description, createdAt)
	return err
}

func (r *ExpenseRepo) GetFiltered(ctx context.Context, filter models.ExpenseFilter) ([]models.Expense, error) {
	var expenses []models.Expense
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", expensesTable)
	args := []interface{}{filter.UserID}
	argIdx := 2

	if filter.Category != "" {
		query += fmt.Sprintf(" AND category = $%d", argIdx)
		args = append(args, filter.Category)
		argIdx++
	}

	if filter.StartDate != "" {
		query += fmt.Sprintf(" AND created_at >= $%d", argIdx)
		args = append(args, filter.StartDate+" 00:00:00")
		argIdx++
	}

	if filter.EndDate != "" {
		query += fmt.Sprintf(" AND created_at <= $%d", argIdx)
		args = append(args, filter.EndDate+" 23:59:59")
		argIdx++
	}

	query += " ORDER BY created_at DESC"

	err := r.db.SelectContext(ctx, &expenses, query, args...)
	if err != nil {
		return nil, err
	}
	return expenses, nil
}
