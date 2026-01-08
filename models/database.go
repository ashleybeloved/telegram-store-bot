package models

import "time"

type User struct {
	ID        uint  `gorm:"primaryKey"`
	UserID    int64 `gorm:"uniqueIndex;not null"`
	Username  string
	Firstname string `gorm:"not null"`
	Lastname  string
	Balance   int64  `gorm:"default:0"`
	LangCode  string `gorm:"column:language_code;default:'ru'"`
	Role      string `gorm:"default:'user'"`
	State     string `gorm:"default:'nothing'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Category struct {
	ID       uint      `gorm:"primaryKey"`
	Name     string    `gorm:"not null"`
	Products []Product `gorm:"constraint:OnDelete:CASCADE;"`
}

type Product struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	CategoryID  uint   `gorm:"not null"`
	Name        string `gorm:"not null"`
	Description string
	Price       int    `gorm:"not null"`
	ImageID     string `gorm:"column:image_id"`
	Stock       int    `gorm:"default:0"`
}
