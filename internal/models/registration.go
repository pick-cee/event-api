package models

import (
	"time"

	"gorm.io/gorm"
)


type Registration struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	EventID   uint           `gorm:"not null" json:"event_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Event     Event          `gorm:"foreignKey:EventID" json:"event,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ensure users cannot register for the same event twice
func (r *Registration) TableName() string {
	return "registrations"
}