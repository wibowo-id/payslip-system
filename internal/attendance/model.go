package attendance

import (
	"time"
)

type Attendance struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	Date      time.Time `gorm:"not null;index:user_date,unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
