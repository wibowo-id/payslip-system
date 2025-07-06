package payslip

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"payslip-system/internal/attendance"
	"payslip-system/internal/auth"
	"payslip-system/internal/overtime"
	"payslip-system/internal/payroll"
	"payslip-system/internal/reimbursement"
	"payslip-system/pkg/middleware"
)

func setupTestRouter(_ *testing.T, db *gorm.DB, userID uint) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.Use(middleware.MockAuthMiddleware(userID))

	h := &Handler{service: NewService(db)}
	api := r.Group("/api")
	api.GET("/payslip/:period_id", h.GeneratePayslip)

	return r
}

func TestGeneratePayslipHandler(t *testing.T) {
	db := setupInMemoryDB(t)
	router := setupTestRouter(t, db, 1)

	req := httptest.NewRequest("GET", "/api/payslip/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), `"period_id":1`)
}

func setupInMemoryDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	err = db.AutoMigrate(
		&auth.User{},
		&attendance.Attendance{},
		&overtime.Overtime{},
		&reimbursement.Reimbursement{},
		&payroll.AttendancePeriod{},
	)
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	// Tambahkan user dummy
	if err := db.Create(&auth.User{
		ID:       1,
		Username: "testuser",
		Password: "hashed",
		Role:     "employee",
	}).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Buat attendance period dummy
	db.Create(&payroll.AttendancePeriod{
		ID:        1,
		StartDate: mustParse("2025-07-01"),
		EndDate:   mustParse("2025-07-05"),
		IsClosed:  true,
	})

	// Simulasi 5 hari kerja
	for i := 1; i <= 5; i++ {
		dateStr := fmt.Sprintf("2025-07-%02d", i)
		date, _ := time.Parse("2006-01-02", dateStr)
		db.Create(&attendance.Attendance{
			UserID: 1,
			Date:   date,
		})
	}

	// Lembur 3 hari x 2 jam
	db.Create(&overtime.Overtime{UserID: 1, Date: mustParse("2025-07-01"), Hours: 2.0})
	db.Create(&overtime.Overtime{UserID: 1, Date: mustParse("2025-07-02"), Hours: 2.0})
	db.Create(&overtime.Overtime{UserID: 1, Date: mustParse("2025-07-03"), Hours: 2.0})

	// Reimbursement 150k
	db.Create(&reimbursement.Reimbursement{UserID: 1, Date: mustParse("2025-07-04"), Amount: 50000})
	db.Create(&reimbursement.Reimbursement{UserID: 1, Date: mustParse("2025-07-05"), Amount: 100000})

	return db
}

func mustParse(dateStr string) time.Time {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		panic(err)
	}
	return t
}
