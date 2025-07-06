package reimbursement

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

func setupTestRouter(db *gorm.DB, userID uint) *gin.Engine {
	r := gin.Default()

	// 🔁 Gunakan mock middleware
	r.Use(middleware.MockAuthMiddleware(userID))

	RegisterRoutes(r.Group("/api"), db)
	return r
}

func TestSubmitReimbursementHandler(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&auth.User{}, &Reimbursement{})

	// Tambah user dummy
	user := auth.User{Username: "john", Password: "hashed", Role: "employee"}
	db.Create(&user)

	r := setupTestRouter(db, user.ID)

	body := SubmitReimbursementRequest{
		Description: "Beli kertas",
		Amount:      50000,
		Date:        time.Date(2025, 7, 3, 0, 0, 0, 0, time.UTC),
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/reimbursements", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}
