package middleware

import (
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func MockAuthMiddleware(userID uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Request.Header.Set("Authorization", "Bearer mocktoken-"+strconv.Itoa(int(userID)))
		c.Next()
	}
}

func GenerateTestJWT(userID uint) string {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret" // fallback
	}
	t, _ := token.SignedString([]byte(secret))
	return t
}
