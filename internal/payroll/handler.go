package payroll

import (
	"net/http"
	"time"

	"payslip-system/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service *Service
}

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	h := &Handler{service: NewService(db)}

	payroll := rg.Group("/payroll", middleware.AuthOnly())
	{
		payroll.POST("/periods", h.CreateAttendancePeriod)
		payroll.POST("/run", h.RunPayroll)
	}
}

func (h *Handler) CreateAttendancePeriod(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only admin can create attendance period"})
		return
	}

	var req AttendancePeriodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.CreateAttendancePeriod(req)
	if err != nil {
		if err == ErrInvalidDateRange {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create period"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "attendance period created"})
}

func (h *Handler) RunPayroll(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only admin can run payroll"})
		return
	}

	var req RunPayrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := uint(c.GetFloat64("user_id"))

	total, err := h.service.RunPayroll(req.AttendancePeriodID, adminID)
	switch err {
	case nil:
		c.JSON(http.StatusCreated, gin.H{
			"message":      "payroll processed successfully",
			"period_id":    req.AttendancePeriodID,
			"total_paid":   total,
			"generated_at": time.Now().Format("2006-01-02 15:04"),
		})
	case ErrPayrollAlreadyRun:
		c.JSON(http.StatusConflict, gin.H{"error": "payroll already run for this period"})
	case ErrPeriodNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "attendance period not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to run payroll"})
	}
}
