package overtime

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

var (
	ErrTooEarly         = errors.New("overtime can only be submitted after working hours")
	ErrTooManyHours     = errors.New("overtime cannot exceed 3 hours")
	ErrAlreadySubmitted = errors.New("overtime already submitted today")
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) SubmitOvertime(userID uint, date time.Time, hours float64) error {
	var existing Overtime
	if err := s.db.Where("user_id = ? AND date = ?", userID, date).First(&existing).Error; err == nil {
		return errors.New("overtime already submitted for this date")
	}

	ot := Overtime{
		UserID: userID,
		Date:   date,
		Hours:  hours,
	}
	return s.db.Create(&ot).Error
}
