package models

import (
	"time"

	"github.com/google/uuid"
	// "gorm.io/gorm"
)

type Account struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	AccountNumber string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"account_number"`
	AccountHolder string    `gorm:"type:varchar(100);not null" json:"account_holder"`
	Balance       int       `gorm:"not null;default:0" json:"balance"`

	BankID uuid.UUID `gorm:"not null" json:"bank_id"`
	Bank   Bank `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"bank"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
