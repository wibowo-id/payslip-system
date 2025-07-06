package reimbursement

import (
	"time"
)

type Reimbursement struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"not null"`
	Amount      float64   `gorm:"not null"`
	Description string    `gorm:"type:text"`
	Date        time.Time `gorm:"not null"` // submission date
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
