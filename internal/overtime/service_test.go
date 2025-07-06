package overtime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&Overtime{})
	return db
}

func TestSubmitOvertime_Success(t *testing.T) {
	db := setupTestDB()
	svc := NewService(db)

	date := time.Date(2025, 7, 5, 0, 0, 0, 0, time.UTC) // Example date
	err := svc.SubmitOvertime(1, date, 2.0)
	assert.Nil(t, err)

	var count int64
	db.Model(&Overtime{}).Where("user_id = ?", 1).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestSubmitOvertime_Duplicate(t *testing.T) {
	db := setupTestDB()
	svc := NewService(db)

	date := time.Date(2025, 7, 5, 0, 0, 0, 0, time.UTC)

	_ = svc.SubmitOvertime(1, date, 1.5)
	err := svc.SubmitOvertime(1, date, 2.0)

	assert.NotNil(t, err)
	assert.Equal(t, "overtime already submitted for this date", err.Error())
}
