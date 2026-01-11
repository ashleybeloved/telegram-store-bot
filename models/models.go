package models

import "time"

type User struct {
	ID        uint  `gorm:"primaryKey"`
	UserID    int64 `gorm:"uniqueIndex;not null"`
	Username  string
	Firstname string `gorm:"not null"`
	Lastname  string
	Balance   int64     `gorm:"default:0"`
	LangCode  string    `gorm:"column:language_code;default:'ru'"`
	State     string    `gorm:"default:'nothing'"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
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
	Description string `gorm:"not null"`
	Price       int64  `gorm:"not null"`
	Stock       int    `gorm:"-"`
	Item        []Item `gorm:"constraint:OnDelete:CASCADE;"`
}

type Item struct {
	ID        uint   `gorm:"primaryKey"`
	ProductID uint   `gorm:"not null;index"`
	Data      string `gorm:"not null"`
	IsSold    bool   `gorm:"default:false"`
}

type PurchasesHistory struct {
	ID          uint   `gorm:"primaryKey"`
	UserID      int64  `gorm:"not null;index"`
	ItemID      uint   `gorm:"not null;index"`
	ProductID   uint   `gorm:"not null;index"`
	Name        string `gorm:"not null"`
	Description string
	Price       int64     `gorm:"not null"`
	Data        string    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

type Promocode struct {
	ID        uint   `gorm:"primaryKey"`
	Code      string `gorm:"uniqueIndex;not null"`
	Reward    int64  `gorm:"not null"`
	MaxUses   int    `gorm:"not null"`
	UsesLeft  int    `gorm:"not null"`
	ExpiresAt time.Time
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type PromocodeUsage struct {
	ID          uint      `gorm:"primaryKey"`
	PromocodeID uint      `gorm:"not null;index"`
	UserID      int64     `gorm:"not null;index"`
	UsedAt      time.Time `gorm:"autoCreateTime"`
}
