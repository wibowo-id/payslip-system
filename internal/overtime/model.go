package overtime

import (
	"time"
)

type Overtime struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_user_date"`
	Date      time.Time `gorm:"not null;uniqueIndex:idx_user_date"`
	Hours     float64   `gorm:"not null"` // 0.5 - 3.0 jam
	CreatedAt time.Time
	UpdatedAt time.Time
}
