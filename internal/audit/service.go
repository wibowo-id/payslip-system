package audit

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

// cara penggunaan log service, ex: submit reimbursement
// auditSvc := audit.NewService(db)
// auditSvc.Log(c, "submit", "reimbursement", reimbursementID)
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Log(c *gin.Context, action, entity string, entityID uint) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	requestID, _ := c.Get("request_id")
	ip := c.ClientIP()

	log := AuditLog{
		UserID:    uint(userID.(float64)), // jwt float64
		Username:  username.(string),
		Action:    action,
		Entity:    entity,
		EntityID:  entityID,
		IP:        ip,
		RequestID: requestID.(string),
	}
	s.db.Create(&log)
}
