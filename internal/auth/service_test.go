package auth

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&User{})
	return db
}

func TestRegisterUser(t *testing.T) {
	db := setupTestDB()
	svc := NewService(db)

	req := RegisterRequest{
		Username: "testuser",
		Password: "secret123",
		Role:     "employee",
	}

	err := svc.RegisterUser(req)
	assert.Nil(t, err)

	var user User
	err = db.Where("username = ?", "testuser").First(&user).Error
	assert.Nil(t, err)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "employee", user.Role)
}

func TestLoginUser(t *testing.T) {
	db := setupTestDB()
	svc := NewService(db)

	// Register first
	_ = svc.RegisterUser(RegisterRequest{
		Username: "testuser",
		Password: "secret123",
		Role:     "employee",
	})

	// Login
	user, err := svc.LoginUser(LoginRequest{
		Username: "testuser",
		Password: "secret123",
	})

	assert.Nil(t, err)
	assert.Equal(t, "testuser", user.Username)
}

func TestGenerateJWT(t *testing.T) {
	db := setupTestDB()
	svc := NewService(db)

	user := &User{
		ID:       1,
		Username: "testuser",
		Role:     "admin",
	}

	os.Setenv("JWT_SECRET", "testsecret")
	token, err := svc.GenerateJWT(user, os.Getenv("JWT_SECRET"))

	assert.Nil(t, err)
	assert.NotEmpty(t, token)
}
