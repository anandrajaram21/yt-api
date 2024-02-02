package models

import (
	"time"

	"gorm.io/gorm"
)

type Video struct {
	gorm.Model             // Adds ID, CreatedAt, UpdatedAt, DeletedAt
	Title        string    `gorm:"size:255;not null"`
	Description  string    `gorm:"type:text;not null"`
	PublishDate  time.Time `gorm:"type:timestamp;index:idx_publish_date"` // Index for sorting by PublishDate
	ThumbnailURL string    `gorm:"size:255"`                              // Thumbnail URL as a string
}
