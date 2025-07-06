package config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB opens a connection to PostgreSQL using GORM
func ConnectDB(cfg *Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Test connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		log.Println("PostgreSQL ping failed:", err)
		return nil, err
	}

	log.Println("Connected to PostgreSQL")
	return db, nil
}
