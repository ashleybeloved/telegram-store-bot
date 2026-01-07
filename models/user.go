package models

import "time"

type User struct {
	ID        int
	UserID    int64
	Username  string
	Firstname string
	Lastname  string
	Balance   int64
	LangCode  string
	Role      string
	State     string
	UpdatedAt time.Time
	CreatedAt time.Time
}
