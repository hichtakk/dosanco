package model

import (
	"time"
)

// Model is common table column
type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `gorm:"created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"deleted_at" json:"deleted_at,omitempty"`
}
