package app

import (
	"payslip-system/internal/attendance"
	"payslip-system/internal/auth"
	"payslip-system/internal/overtime"
	"payslip-system/internal/payroll"
	"payslip-system/internal/payslip"
	"payslip-system/internal/reimbursement"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupRouterForTest returns a gin.Engine for integration testing
func SetupRouterForTest() *gin.Engine {
	// Use SQLite in-memory for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to open test database: " + err.Error())
	}

	// Auto-migrate all relevant models
	if err := db.AutoMigrate(
		&auth.User{},
		&attendance.Attendance{},
		&overtime.Overtime{},
		&reimbursement.Reimbursement{},
		&payroll.AttendancePeriod{},
		&payroll.Payroll{},
	); err != nil && !strings.Contains(err.Error(), "already exists") {
		panic("failed to migrate test db: " + err.Error())
	}

	// Optionally seed test users (if needed)
	if err := auth.SeedUsers(db); err != nil {
		panic("failed to seed users in test: " + err.Error())
	}

	// Create Gin engine
	r := gin.Default()

	// Register all routes under /api
	api := r.Group("/api")
	auth.RegisterRoutes(api, db)
	attendance.RegisterRoutes(api, db, false)
	overtime.RegisterRoutes(api, db)
	reimbursement.RegisterRoutes(api, db)
	payslip.RegisterRoutes(api, db)
	payroll.RegisterRoutes(api, db)

	return r
}
