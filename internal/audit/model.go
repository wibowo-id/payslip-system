package audit

import (
	"time"
)

type AuditLog struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	Username  string
	Action    string
	Entity    string
	EntityID  uint
	IP        string
	RequestID string
	CreatedAt time.Time
}
