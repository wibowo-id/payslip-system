package payslip

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

type Payslip struct {
	PeriodID       uint                  `json:"period_id"`
	PeriodStart    string                `json:"period_start"`
	PeriodEnd      string                `json:"period_end"`
	Attendance     int64                 `json:"attendance_days"`
	OvertimeHours  float64               `json:"overtime_hours"`
	Reimbursements []ReimbursementDetail `json:"reimbursements"`
	SalaryDaily    float64               `json:"salary_daily"`
	SalaryPart     float64               `json:"salary_component"`
	OvertimePart   float64               `json:"overtime_component"`
	ReimbursePart  float64               `json:"reimburse_component"`
	TotalTakeHome  float64               `json:"total_take_home"`
}

type ReimbursementDetail struct {
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

var ErrNotFound = errors.New("period not found or not processed")

func (s *Service) GeneratePayslip(userID uint, periodID uint) (*Payslip, error) {
	// Load period
	var period struct {
		ID        uint
		StartDate time.Time
		EndDate   time.Time
		IsClosed  bool
	}
	err := s.db.Raw(`
		SELECT id, start_date, end_date, is_closed 
		FROM attendance_periods 
		WHERE id = ?`, periodID).Scan(&period).Error
	if err != nil || !period.IsClosed {
		return nil, ErrNotFound
	}

	// Load salary
	var salary float64
	s.db.Raw(`SELECT salary FROM users WHERE id = ?`, userID).Scan(&salary)

	daily := salary / 20

	// Count attendance
	var attendCount int64
	s.db.Raw(`SELECT COUNT(*) FROM attendances 
		WHERE user_id = ? AND date BETWEEN ? AND ?`,
		userID, period.StartDate, period.EndDate).Scan(&attendCount)

	// Sum overtime
	var otTotal float64
	s.db.Raw(`SELECT SUM(hours) FROM overtimes 
		WHERE user_id = ? AND date BETWEEN ? AND ?`,
		userID, period.StartDate, period.EndDate).Scan(&otTotal)

	// List reimbursements
	var reimbursements []ReimbursementDetail
	s.db.Raw(`SELECT amount, description, date FROM reimbursements 
		WHERE user_id = ? AND date BETWEEN ? AND ? 
		ORDER BY date`,
		userID, period.StartDate, period.EndDate).Scan(&reimbursements)

	var reimburseTotal float64
	for _, r := range reimbursements {
		reimburseTotal += r.Amount
	}

	salaryPart := float64(attendCount) * daily
	overtimePart := otTotal * daily * 2
	total := salaryPart + overtimePart + reimburseTotal

	return &Payslip{
		PeriodID:       period.ID,
		PeriodStart:    period.StartDate.Format("2006-01-02"),
		PeriodEnd:      period.EndDate.Format("2006-01-02"),
		Attendance:     attendCount,
		OvertimeHours:  otTotal,
		Reimbursements: reimbursements,
		SalaryDaily:    daily,
		SalaryPart:     salaryPart,
		OvertimePart:   overtimePart,
		ReimbursePart:  reimburseTotal,
		TotalTakeHome:  total,
	}, nil
}

type PayslipSummary struct {
	PeriodID    uint              `json:"period_id"`
	PeriodStart string            `json:"period_start"`
	PeriodEnd   string            `json:"period_end"`
	Employees   []EmployeePayslip `json:"employees"`
	TotalPaid   float64           `json:"total_paid"`
}

type EmployeePayslip struct {
	UserID   uint    `json:"user_id"`
	Name     string  `json:"name"`
	TakeHome float64 `json:"take_home"`
}

func (s *Service) GenerateSummary(periodID uint) (*PayslipSummary, error) {
	// Load attendance period
	var period struct {
		ID        uint
		StartDate time.Time
		EndDate   time.Time
		IsClosed  bool
	}
	err := s.db.Raw(`
		SELECT id, start_date, end_date, is_closed 
		FROM attendance_periods 
		WHERE id = ?`, periodID).Scan(&period).Error
	if err != nil || !period.IsClosed {
		return nil, ErrNotFound
	}

	// Get employees
	var users []struct {
		ID     uint
		Name   string
		Salary float64
	}
	err = s.db.Raw(`
		SELECT id, name, salary FROM users WHERE role = 'employee'`).Scan(&users).Error
	if err != nil {
		return nil, err
	}

	var employees []EmployeePayslip
	var grandTotal float64

	for _, u := range users {
		daily := u.Salary / 20

		var attend int64
		s.db.Raw(`SELECT COUNT(*) FROM attendances 
			WHERE user_id = ? AND date BETWEEN ? AND ?`,
			u.ID, period.StartDate, period.EndDate).Scan(&attend)

		var ot float64
		s.db.Raw(`SELECT SUM(hours) FROM overtimes 
			WHERE user_id = ? AND date BETWEEN ? AND ?`,
			u.ID, period.StartDate, period.EndDate).Scan(&ot)

		var reimburse float64
		s.db.Raw(`SELECT SUM(amount) FROM reimbursements 
			WHERE user_id = ? AND date BETWEEN ? AND ?`,
			u.ID, period.StartDate, period.EndDate).Scan(&reimburse)

		takeHome := float64(attend)*daily + ot*daily*2 + reimburse
		grandTotal += takeHome

		employees = append(employees, EmployeePayslip{
			UserID:   u.ID,
			Name:     u.Name,
			TakeHome: takeHome,
		})
	}

	return &PayslipSummary{
		PeriodID:    period.ID,
		PeriodStart: period.StartDate.Format("2006-01-02"),
		PeriodEnd:   period.EndDate.Format("2006-01-02"),
		Employees:   employees,
		TotalPaid:   grandTotal,
	}, nil
}
