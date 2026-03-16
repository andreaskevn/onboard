package models

import "time"

type Item struct {
	ID          int       `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}