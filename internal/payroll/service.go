package payroll

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

var (
	ErrPayrollAlreadyRun = errors.New("payroll already processed for this period")
	ErrPeriodNotFound    = errors.New("attendance period not found")
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateAttendancePeriod(req AttendancePeriodRequest) error {
	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return err
	}
	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return err
	}
	if end.Before(start) {
		return ErrInvalidDateRange
	}

	period := AttendancePeriod{
		StartDate: start,
		EndDate:   end,
	}
	return s.db.Create(&period).Error
}

func (s *Service) RunPayroll(periodID uint, adminID uint) (float64, error) {
	var period AttendancePeriod
	if err := s.db.First(&period, periodID).Error; err != nil {
		return 0, ErrPeriodNotFound
	}

	// Check if payroll already run
	var existing Payroll
	if err := s.db.Where("attendance_period_id = ?", periodID).First(&existing).Error; err == nil {
		return 0, ErrPayrollAlreadyRun
	}

	// Get employee list
	var users []struct {
		ID     uint
		Salary float64
	}
	err := s.db.Raw(`
		SELECT id, salary FROM users WHERE role = 'employee'
	`).Scan(&users).Error
	if err != nil {
		return 0, err
	}

	total := 0.0

	// Loop each employee, calculate payslip
	for _, u := range users {
		var attendCount int64
		s.db.Raw(`
			SELECT COUNT(*) FROM attendances 
			WHERE user_id = ? AND date BETWEEN ? AND ?`,
			u.ID, period.StartDate, period.EndDate).Scan(&attendCount)

		var otTotal float64
		s.db.Raw(`
			SELECT SUM(hours) FROM overtimes 
			WHERE user_id = ? AND date BETWEEN ? AND ?`,
			u.ID, period.StartDate, period.EndDate).Scan(&otTotal)

		var reimburseTotal float64
		s.db.Raw(`
			SELECT SUM(amount) FROM reimbursements 
			WHERE user_id = ? AND date BETWEEN ? AND ?`,
			u.ID, period.StartDate, period.EndDate).Scan(&reimburseTotal)

		// assume 20 work days
		daily := u.Salary / 20
		salaryComponent := float64(attendCount) * daily
		overtimeComponent := otTotal * daily * 2
		takeHome := salaryComponent + overtimeComponent + reimburseTotal

		total += takeHome

		// optional: save detail per employee if kamu buat `PayslipDetail` table
	}

	// Mark period as closed
	s.db.Model(&period).Update("is_closed", true)

	// Save payroll record
	payroll := Payroll{
		AttendancePeriodID: periodID,
		GeneratedAt:        time.Now(),
		CreatedBy:          adminID,
		TotalPaid:          total,
	}
	if err := s.db.Create(&payroll).Error; err != nil {
		return 0, err
	}

	return total, nil
}
