package app

import (
	"payslip-system/config"
	"payslip-system/internal/attendance"
	"payslip-system/internal/auth"
	"payslip-system/internal/overtime"
	"payslip-system/internal/payroll"
	"payslip-system/internal/payslip"
	"payslip-system/internal/reimbursement"
	"payslip-system/pkg/logger"
	"payslip-system/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupApp() (*gin.Engine, *gorm.DB, error) {
	// Load config (DB, env, dsb)
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, nil, err
	}

	// Connect DB
	db, err := config.ConnectDB(cfg)
	if err != nil {
		return nil, nil, err
	}

	// Setup logger
	logger.InitLogger()

	// Init Gin
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger())       // logrus or zap
	r.Use(middleware.RequestIDMiddleware()) // inject X-Request-ID

	// Route group
	api := r.Group("/api")

	// Register routes
	auth.RegisterRoutes(api, db)
	attendance.RegisterRoutes(api, db, true)
	overtime.RegisterRoutes(api, db)
	reimbursement.RegisterRoutes(api, db)
	payroll.RegisterRoutes(api, db)
	payslip.RegisterRoutes(api, db)

	return r, db, nil
}
