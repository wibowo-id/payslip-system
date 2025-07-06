package attendance

import (
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

func setupRouterWithAuth(t *testing.T, mockNow time.Time) (*gin.Engine, *gorm.DB, auth.User) {
	gin.SetMode(gin.TestMode)

	// Setup in-memory DB
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&auth.User{}, &Attendance{})
	assert.NoError(t, err)

	// Dummy user
	user := auth.User{
		Username: "john",
		Password: "$2a$10$h8muj/0DYvXq1KfNaRLpEO9hvogvuhMNpzkUHaGljhf6Lq.CpAfoa",
		Role:     "employee",
	}
	db.Create(&user)

	service := NewService(db)
	service.NowFunc = func() time.Time {
		return mockNow
	}
	handler := &Handler{service: service}

	// Setup router with mocked auth middleware
	r := gin.Default()
	r.Use(middleware.MockAuthMiddleware(user.ID))
	RegisterRoutes(r.Group("/api"), handler, false)
	r.POST("/attendance/checkin", handler.CheckIn)

	return r, db, user
}

func TestCheckInHandler_Success(t *testing.T) {
	// Senin (bukan weekend)
	mockTime := time.Date(2025, 7, 7, 9, 0, 0, 0, time.Local)
	router, _, _ := setupRouterWithAuth(t, mockTime)

	req := httptest.NewRequest(http.MethodPost, "/api/attendance/checkin", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCheckInHandler_Weekend(t *testing.T) {
	// Minggu
	mockTime := time.Date(2025, 7, 6, 9, 0, 0, 0, time.Local)
	router, _, _ := setupRouterWithAuth(t, mockTime)

	req := httptest.NewRequest(http.MethodPost, "/api/attendance/checkin", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
