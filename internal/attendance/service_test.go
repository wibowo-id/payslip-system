package attendance

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"payslip-system/internal/auth"
)

func setupTestService(t *testing.T, now time.Time) (*Service, *gorm.DB, auth.User) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	assert.NoError(t, db.AutoMigrate(&Attendance{}, &auth.User{}))

	user := auth.User{
		Username: "testuser",
		Password: "dummy",
		Role:     "employee",
	}
	db.Create(&user)

	service := NewService(db)
	service.NowFunc = func() time.Time {
		return now
	}

	return service, db, user
}

func TestSubmitAttendance_Success(t *testing.T) {
	// Monday (weekday)
	now := time.Date(2025, 7, 7, 9, 0, 0, 0, time.Local)
	service, _, user := setupTestService(t, now)

	err := service.SubmitAttendance(user.ID)
	assert.NoError(t, err)

	var count int64
	service.db.Model(&Attendance{}).Where("user_id = ?", user.ID).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestSubmitAttendance_Weekend(t *testing.T) {
	// Sunday
	now := time.Date(2025, 7, 6, 9, 0, 0, 0, time.Local)
	service, _, user := setupTestService(t, now)

	err := service.SubmitAttendance(user.ID)
	assert.Equal(t, ErrWeekend, err)
}
