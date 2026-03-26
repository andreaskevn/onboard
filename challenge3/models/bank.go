package models

import (
	"time"

	"github.com/google/uuid"
)

type Bank struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"not null"`
	Code      string    `gorm:"unique;not null"`
	Accounts  []Account `gorm:"foreignKey:BankID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
