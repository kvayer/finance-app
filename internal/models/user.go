package models

import "time"

type User struct {
	ID           int       `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}

type RegisterRequest struct {
	Name            string
	Email           string
	Password        string
	ConfirmPassword string
}

type LoginRequest struct {
	Email    string
	Password string
}

type UpdatePasswordRequest struct {
	OldPassword string
	NewPassword string
}
