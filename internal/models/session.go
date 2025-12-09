package models

import "time"

type Session struct {
	Token  string    `db:"token"`
	UserID int       `db:"user_id"`
	Expiry time.Time `db:"expiry"`
}
