package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "challenge3/models"
)

func InitDb() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=End.Game2o2o dbname=challenge3 port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("Gagal connect DB: %w", err)
	}

	return db, nil
}
