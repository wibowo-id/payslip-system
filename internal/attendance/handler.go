package attendance

import (
	"net/http"

	"payslip-system/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, useAuth bool) {
	attendance := rg.Group("/attendance")
	if useAuth {
		attendance.Use(middleware.AuthOnly())
	}
	attendance.POST("/checkin", h.CheckIn)
}

func (h *Handler) CheckIn(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var userID uint
	switch v := userIDVal.(type) {
	case float64:
		userID = uint(v)
	case uint:
		userID = v
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID type"})
		return
	}

	err := h.service.SubmitAttendance(userID)
	switch err {
	case nil:
		c.JSON(http.StatusCreated, AttendanceResponse{
			Message: "attendance recorded",
			Date:    h.service.NowFunc().Format("2006-01-02"),
		})
	case ErrWeekend:
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot check in on weekend"})
	case ErrAlreadyCheckedIn:
		c.JSON(http.StatusConflict, gin.H{"error": "already checked in today"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record attendance"})
	}
}

func AuthOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := c.Get("user_id"); !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
