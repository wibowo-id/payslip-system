package overtime

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

	overtime := rg.Group("/overtime", middleware.AuthOnly())
	overtime.POST("", h.SubmitOvertime)
}

func (h *Handler) SubmitOvertime(c *gin.Context) {
	userID := uint(c.GetFloat64("user_id"))

	var req OvertimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.SubmitOvertime(userID, time.Now(), req.Hours)
	switch err {
	case nil:
		c.JSON(http.StatusCreated, gin.H{
			"message": "overtime submitted",
			"hours":   req.Hours,
			"date":    time.Now().Format("2006-01-02"),
		})
	case ErrTooEarly:
		c.JSON(http.StatusForbidden, gin.H{"error": "overtime can only be submitted after 17:00"})
	case ErrTooManyHours:
		c.JSON(http.StatusBadRequest, gin.H{"error": "overtime max 3 hours per day"})
	case ErrAlreadySubmitted:
		c.JSON(http.StatusConflict, gin.H{"error": "already submitted overtime today"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to submit overtime"})
	}
}
