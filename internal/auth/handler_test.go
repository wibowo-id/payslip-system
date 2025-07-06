package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&User{})

	RegisterRoutes(r.Group("/api"), db)
	return r
}

func TestRegisterHandler(t *testing.T) {
	r := setupTestRouter()

	payload := RegisterRequest{
		Username: "handleruser",
		Password: "password",
		Role:     "employee",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestLoginHandler(t *testing.T) {
	r := setupTestRouter()

	// Register user dulu
	registerPayload := RegisterRequest{
		Username: "loginuser",
		Password: "mypassword",
		Role:     "employee",
	}
	body, _ := json.Marshal(registerPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Set env for JWT
	os.Setenv("JWT_SECRET", "mytestsecret")

	// Login
	loginPayload := LoginRequest{
		Username: "loginuser",
		Password: "mypassword",
	}
	body, _ = json.Marshal(loginPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.NotEmpty(t, resp["token"])
	assert.Equal(t, "loginuser", resp["username"])
}
