package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex" json:"username"`
	Email     string         `gorm:"uniqueIndex" json:"email"`
	Password  string         `json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Email represents a registered email address
type Email struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	UserID    string         `gorm:"index" json:"user_id"`
	Address   string         `gorm:"uniqueIndex" json:"address"`
	Verified  bool           `json:"verified"`
	Token     string         `json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Alias represents an email alias
type Alias struct {
	ID           string         `gorm:"primaryKey" json:"id"`
	UserID       string         `gorm:"index" json:"user_id"`
	EmailID      string         `gorm:"index" json:"email_id"`
	AliasAddress string         `gorm:"uniqueIndex" json:"alias_address"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
