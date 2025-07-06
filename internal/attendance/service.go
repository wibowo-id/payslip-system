package attendance

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

var ErrWeekend = errors.New("cannot check in on weekend")
var ErrAlreadyCheckedIn = errors.New("already checked in today")

type Service struct {
	db      *gorm.DB
	NowFunc func() time.Time
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db:      db,
		NowFunc: time.Now,
	}
}

func (s *Service) SubmitAttendance(userID uint) error {
	now := s.NowFunc()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	weekday := today.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return ErrWeekend
	}

	var existing Attendance
	if err := s.db.Where("user_id = ? AND date = ?", userID, today).First(&existing).Error; err == nil {
		return ErrAlreadyCheckedIn
	}

	record := Attendance{
		UserID: userID,
		Date:   today,
	}
	return s.db.Create(&record).Error
}
