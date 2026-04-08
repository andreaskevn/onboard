package models

import (
	"time"
)

type IdempotencyKey struct {
	ID        string `gorm:"primaryKey"` // idempotency key
	Response  string `gorm:"type:text"`  // simpan response JSON
	Status    int    // HTTP status
	CreatedAt time.Time
}
