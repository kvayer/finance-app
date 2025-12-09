package models

import "time"

type Expense struct {
	ID          int       `db:"id"`
	UserID      int       `db:"user_id"`
	Category    string    `db:"category"`
	Amount      float64   `db:"amount"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}

type ExpenseFilter struct {
	UserID    int
	Category  string
	StartDate string
	EndDate   string
}
