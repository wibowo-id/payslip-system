package overtime

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"payslip-system/internal/auth"
	"payslip-system/pkg/middleware"
)

type SubmitOvertimeRequest struct {
	Date  time.Time `json:"date"`
	Hours float64   `json:"hours" binding:"required,gt=0,lte=3"`
}

func setupOvertimeRouter(db *gorm.DB, userID uint) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// 🔁 Gunakan mock middleware
	r.Use(middleware.MockAuthMiddleware(userID))

	RegisterRoutes(r.Group("/api"), db)
	return r
}

func TestSubmitOvertimeHandler_Success(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&auth.User{}, &Overtime{})

	user := auth.User{Username: "john", Password: "hashed", Role: "employee"}
	db.Create(&user)

	r := setupOvertimeRouter(db, user.ID)

	payload := SubmitOvertimeRequest{
		Date:  time.Date(2025, 7, 5, 0, 0, 0, 0, time.UTC),
		Hours: 2.0,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/overtime", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestSubmitOvertimeHandler_Duplicate(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&auth.User{}, &Overtime{})

	user := auth.User{Username: "john", Password: "hashed", Role: "employee"}
	db.Create(&user)

	r := setupOvertimeRouter(db, user.ID)

	payload := SubmitOvertimeRequest{
		Date:  time.Date(2025, 7, 5, 0, 0, 0, 0, time.UTC),
		Hours: 2.0,
	}
	body, _ := json.Marshal(payload)

	// Submit pertama
	req := httptest.NewRequest(http.MethodPost, "/api/overtime", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Submit kedua (duplikat)
	req2 := httptest.NewRequest(http.MethodPost, "/api/overtime", bytes.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusBadRequest, w2.Code)
}
