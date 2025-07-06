package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) error {
	// Only seed if empty
	var count int64
	db.Model(&User{}).Count(&count)
	if count > 0 {
		return nil
	}

	password, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	// Seed 1 admin
	admin := User{
		Username: "admin",
		Password: string(password),
		Role:     "admin",
	}
	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	// Seed 100 employees
	var users []User
	for i := 1; i <= 100; i++ {
		users = append(users, User{
			Username: fmt.Sprintf("employee%d", i),
			Password: string(password),
			Role:     "employee",
		})
	}
	return db.Create(&users).Error
}
