package reimbursement

import (
	"time"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Submit(userID uint, req ReimbursementRequest) error {
	record := Reimbursement{
		UserID:      userID,
		Amount:      req.Amount,
		Description: req.Description,
		Date:        time.Now(),
	}
	return s.db.Create(&record).Error
}

func (s *Service) SubmitReimbursement(userID uint, description string, amount float64, date time.Time) error {
	reimb := Reimbursement{
		UserID:      userID,
		Description: description,
		Amount:      amount,
		Date:        date,
	}
	return s.db.Create(&reimb).Error
}
