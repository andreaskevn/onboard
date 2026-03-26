package models

import (
	"time"

	"github.com/google/uuid"
	// "gorm.io/gorm"
)

type Transaction struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	AccountFrom uuid.UUID `gorm:"type:uuid;not null" json:"account_from"`
	AccountTo   uuid.UUID `gorm:"type:uuid;not null" json:"account_to"`
	Amount      int       `gorm:"type:int;not null" json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
	AdminFee    int       `gorm:"not null;default:0" json:"admin_fee"`
}
