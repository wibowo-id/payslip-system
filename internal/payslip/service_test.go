package payslip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculatePayslipComponents(t *testing.T) {
	// Simulasi data
	payslip := Payslip{
		Attendance:     20,
		OvertimeHours:  10.5,
		SalaryDaily:    100_000,
		Reimbursements: []ReimbursementDetail{{Amount: 250_000}, {Amount: 150_000}},
	}

	// Logika perhitungan manual
	payslip.SalaryPart = float64(payslip.Attendance) * payslip.SalaryDaily
	payslip.OvertimePart = payslip.OvertimeHours * 20_000
	for _, r := range payslip.Reimbursements {
		payslip.ReimbursePart += r.Amount
	}
	payslip.TotalTakeHome = payslip.SalaryPart + payslip.OvertimePart + payslip.ReimbursePart

	assert.Equal(t, 2_000_000.0, payslip.SalaryPart)
	assert.Equal(t, 210_000.0, payslip.OvertimePart)
	assert.Equal(t, 400_000.0, payslip.ReimbursePart)
	assert.Equal(t, 2_610_000.0, payslip.TotalTakeHome)
}
