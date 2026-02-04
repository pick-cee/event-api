package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Title         string         `gorm:"not null" json:"title"`
	Description   string         `json:"description"`
	Location      string         `gorm:"not null" json:"location"`
	DateTime      time.Time      `gorm:"not null" json:"date_time"`
	CreatorID     uint           `gorm:"not null" json:"creator_id"`
	Creator       User           `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	Registrations []Registration `gorm:"foreignKey:EventID" json:"registrations,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
