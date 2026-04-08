package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "challenge3/models"
	// "github.com/joho/godotenv"
	"os"
)

func InitDb() (*gorm.DB, error) {
	// err := godotenv.Load("../.env")
	// if err != nil {
	// 	return nil, fmt.Errorf("error loading .env file: %w", err)
	// }
	// dsn := os.Getenv("DB_URL")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	fmt.Println("DSN:", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
