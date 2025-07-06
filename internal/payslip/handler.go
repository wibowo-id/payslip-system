package payslip

import (
	"net/http"
	"strconv"

	"payslip-system/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service *Service
}

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	h := &Handler{service: NewService(db)}
	payslip := rg.Group("payslip", middleware.AuthOnly())
	{
		payslip.GET("/:period_id", h.GeneratePayslip)
		payslip.GET("/summary/:period_id", h.GenerateSummary)
	}
}

func (h *Handler) GeneratePayslip(c *gin.Context) {
	userID := uint(c.GetFloat64("user_id"))

	periodIDStr := c.Param("period_id")
	periodID, err := strconv.Atoi(periodIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period_id"})
		return
	}

	result, err := h.service.GeneratePayslip(userID, uint(periodID))
	switch err {
	case nil:
		c.JSON(http.StatusOK, result)
	case ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "attendance period not found or payroll not yet run"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate payslip"})
	}
}

func (h *Handler) GenerateSummary(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only admin can view payslip summary"})
		return
	}

	periodIDStr := c.Param("period_id")
	periodID, err := strconv.Atoi(periodIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period_id"})
		return
	}

	result, err := h.service.GenerateSummary(uint(periodID))
	switch err {
	case nil:
		c.JSON(http.StatusOK, result)
	case ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "period not found or payroll not yet run"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate summary"})
	}
}
