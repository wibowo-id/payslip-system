package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Role     string `gorm:"not null"` // e.g. admin, employee
}
