package models

import (
	"time"

	"github.com/google/uuid"
)

type CalendarEvent struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CoupleID    uuid.UUID  `gorm:"type:uuid;not null" json:"couple_id"`
	CreatedBy   uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	Title       string     `gorm:"type:varchar(150);not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	EventDate   time.Time  `gorm:"not null" json:"event_date"`
	EventType   string     `gorm:"type:varchar(50);not null" json:"event_type"` // custom, memory, date_plan, anniversary
	ReferenceID *uuid.UUID `gorm:"type:uuid" json:"reference_id"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relasi
	Couple Couple `gorm:"foreignKey:CoupleID" json:"-"`
	Author User   `gorm:"foreignKey:CreatedBy" json:"-"`
}