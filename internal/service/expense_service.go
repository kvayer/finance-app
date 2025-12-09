package service

import (
	"context"
	"finance-tracker/internal/models"
	"finance-tracker/internal/repository"
	"time"
)

type ExpenseService struct {
	repo *repository.ExpenseRepo
}

func NewExpenseService(repo *repository.ExpenseRepo) *ExpenseService {
	return &ExpenseService{repo: repo}
}

func (s *ExpenseService) AddExpense(ctx context.Context, userID int, category string, amount float64, desc string, dateStr string) error {
	createdAt := time.Now()
	// Parse da data do input HTML type="datetime-local"
	if dateStr != "" {
		if t, err := time.Parse("2006-01-02T15:04", dateStr); err == nil {
			createdAt = t
		}
	}

	e := &models.Expense{
		UserID:      userID,
		Category:    category,
		Amount:      amount,
		Description: desc,
		CreatedAt:   createdAt,
	}
	return s.repo.Create(ctx, e)
}

func (s *ExpenseService) GetFilteredExpenses(ctx context.Context, filter models.ExpenseFilter) ([]models.Expense, error) {
	return s.repo.GetFiltered(ctx, filter)
}

func (s *ExpenseService) CalculateTotal(expenses []models.Expense) float64 {
	var total float64
	for _, e := range expenses {
		total += e.Amount
	}
	return total
}
