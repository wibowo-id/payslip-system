package payroll

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"payslip-system/internal/auth"
	"payslip-system/pkg/middleware"
)

func setupTestRouter(t *testing.T, userID uint) (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&auth.User{}, &AttendancePeriod{})
	assert.NoError(t, err)

	r := gin.Default()

	// 🔁 Gunakan MockAuthMiddleware
	r.Use(middleware.MockAuthMiddleware(userID))

	rg := r.Group("/api")
	RegisterRoutes(rg, db)

	return r, db
}

func TestCreateAttendancePeriodHandler(t *testing.T) {
	// Buat dummy user
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&auth.User{})
	user := auth.User{Username: "admin", Password: "hashed", Role: "admin"}
	db.Create(&user)

	router, _ := setupTestRouter(t, user.ID)

	body := map[string]string{
		"start_date": "2025-07-01",
		"end_date":   "2025-07-31",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/api/payroll/periods", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}
