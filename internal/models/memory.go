package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Memory struct {
	ID                      uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CoupleID                uuid.UUID      `gorm:"type:uuid;not null" json:"couple_id"`
	CreatedBy               uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	Title                   string         `gorm:"type:varchar(150);not null" json:"title"`
	Description             string         `gorm:"type:text" json:"description"`
	MemoryDate              time.Time      `gorm:"type:date;not null" json:"memory_date"`
	LocationName            string         `gorm:"type:varchar(150)" json:"location_name"`
	Latitude                *float64       `gorm:"type:decimal(10,8)" json:"latitude"`
	Longitude               *float64       `gorm:"type:decimal(11,8)" json:"longitude"`
	Mood                    string         `gorm:"type:varchar(50)" json:"mood"`
	Tags                    pq.StringArray `gorm:"type:text[]" json:"tags"`
	ConvertedFromDatePlanID *uuid.UUID     `gorm:"type:uuid" json:"-"`
	CreatedAt               time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt               time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt               gorm.DeletedAt `gorm:"index" json:"-"`

	Couple Couple `gorm:"foreignKey:CoupleID" json:"-"`
	Author User   `gorm:"foreignKey:CreatedBy" json:"-"`
}