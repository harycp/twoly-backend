package models

import (
	"time"

	"github.com/google/uuid"
)

type DatePlan struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CoupleID     uuid.UUID `gorm:"type:uuid;not null" json:"couple_id"`
	CreatedBy    uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	Title        string    `gorm:"type:varchar(150);not null" json:"title"`
	Notes        string    `gorm:"type:text" json:"notes"`
	PlanDate     time.Time `gorm:"not null" json:"plan_date"`
	LocationName string    `gorm:"type:varchar(150)" json:"location_name"`
	Latitude     *float64  `gorm:"type:decimal(10,8)" json:"latitude"`
	Longitude    *float64  `gorm:"type:decimal(11,8)" json:"longitude"`
	Status       string    `gorm:"type:varchar(20);default:'planned'" json:"status"` // planned, ongoing, completed, cancelled
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	Checklists []DatePlanChecklist `gorm:"foreignKey:DatePlanID;constraint:OnDelete:CASCADE;" json:"checklists"`
	Couple     Couple              `gorm:"foreignKey:CoupleID" json:"-"`
	Author     User                `gorm:"foreignKey:CreatedBy" json:"-"`
}

type DatePlanChecklist struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	DatePlanID uuid.UUID `gorm:"type:uuid;not null" json:"date_plan_id"`
	Item       string    `gorm:"type:text;not null" json:"item"`
	IsChecked  bool      `gorm:"default:false" json:"is_checked"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}