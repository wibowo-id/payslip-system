package reimbursement

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

type SubmitReimbursementRequest struct {
	Description string    `json:"description" binding:"required"`
	Amount      float64   `json:"amount" binding:"required,gt=0"`
	Date        time.Time `json:"date" binding:"required"` // ISO format
}

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	h := &Handler{service: NewService(db)}

	routes := rg.Group("/reimbursements", middleware.AuthOnly())
	routes.POST("", h.SubmitReimbursement)
}

func (h *Handler) SubmitReimbursement(c *gin.Context) {
	userID := uint(c.GetFloat64("user_id"))

	var req ReimbursementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Submit(userID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to submit reimbursement"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "reimbursement submitted",
		"amount":      req.Amount,
		"description": req.Description,
		"date":        time.Now().Format("2006-01-02"),
	})
}
