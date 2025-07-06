package payroll

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&AttendancePeriod{})
	assert.NoError(t, err)
	return db
}

func TestCreateAttendancePeriod(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db)

	startDate := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 7, 31, 0, 0, 0, 0, time.UTC)

	req := AttendancePeriodRequest{
		StartDate: startDate.Format("2006-01-02"),
		EndDate:   endDate.Format("2006-01-02"),
	}

	err := service.CreateAttendancePeriod(req)
	assert.NoError(t, err)

	var period AttendancePeriod
	err = db.First(&period).Error
	assert.NoError(t, err)
	assert.Equal(t, startDate, period.StartDate)
	assert.Equal(t, endDate, period.EndDate)
}
