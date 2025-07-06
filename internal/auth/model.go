package auth

import (
	"time"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"unique;not null"`
	Password  string `gorm:"not null"` // hashed password
	Role      string `gorm:"not null"` // "admin" or "employee"
	CreatedAt time.Time
	UpdatedAt time.Time
}
