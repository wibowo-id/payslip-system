package main

import (
	"log"

	"payslip-system/config"
	"payslip-system/internal/app"
	"payslip-system/internal/attendance"
	"payslip-system/internal/audit"
	"payslip-system/internal/auth"
	"payslip-system/internal/overtime"
	"payslip-system/internal/payroll"
	"payslip-system/internal/reimbursement"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Connect DB
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// AutoMigrate
	err = db.AutoMigrate(
		&auth.User{},
		&audit.AuditLog{},
		&payroll.AttendancePeriod{},
		&payroll.Payroll{},
		&attendance.Attendance{},
		&overtime.Overtime{},
		&reimbursement.Reimbursement{},
	)
	if err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	// Seed users
	if err := auth.SeedUsers(db); err != nil {
		log.Fatalf("failed to seed users: %v", err)
	}

	// Setup Gin + routes via app
	r, _, err := app.SetupApp()
	if err != nil {
		log.Fatalf("failed to setup app: %v", err)
	}

	// Start server
	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
