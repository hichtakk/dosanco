package model

import (
	"time"
)

type Model struct {
	ID			uint		`gorm:"primary_key" json:"id"`
	CreatedAt	time.Time	`gorm:"created_at" json:"created_at"`
	UpdatedAt	time.Time	`gorm:"created_at" json:"updated_at"`
	DeletedAt	*time.Time	`gorm:"created_at" json:"deleted_at,omitempty"`
}