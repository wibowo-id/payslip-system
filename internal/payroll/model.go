package payroll

import (
	"time"
)

type AttendancePeriod struct {
	ID        uint      `gorm:"primaryKey"`
	StartDate time.Time `gorm:"not null"`
	EndDate   time.Time `gorm:"not null"`
	IsClosed  bool      `gorm:"default:false"` // true if payroll has been run
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Payroll struct {
	ID                 uint      `gorm:"primaryKey"`
	AttendancePeriodID uint      `gorm:"not null;unique"` // only once per period
	GeneratedAt        time.Time `gorm:"not null"`
	CreatedBy          uint      `gorm:"not null"` // admin user
	TotalPaid          float64   `gorm:"not null"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
