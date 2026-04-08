package models

import (
	"time"

	"github.com/google/uuid"
)

type Bank struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Code      string    `gorm:"unique;not null" json:"code"`
	Accounts  []Account `gorm:"foreignKey:BankID" json:"accounts"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
