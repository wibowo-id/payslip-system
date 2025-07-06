package auth

import (
	"errors"

	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) RegisterUser(req RegisterRequest) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := User{
		Username: req.Username,
		Password: string(hashed),
		Role:     req.Role,
	}
	return s.db.Create(&user).Error
}

func (s *Service) LoginUser(req LoginRequest) (*User, error) {
	var user User
	if err := s.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return nil, errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	return &user, nil
}

func (s *Service) GenerateJWT(user *User, jwtSecret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
