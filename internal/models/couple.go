package models

import (
	"time"

	"github.com/google/uuid"
)

type Couple struct {
	ID              uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserOneID       uuid.UUID  `gorm:"type:uuid;not null" json:"user_one_id"`
	UserTwoID       *uuid.UUID `gorm:"type:uuid" json:"user_two_id"` // Bisa null jika pasangan belum join
	InviteCode      string     `gorm:"type:varchar(20);unique" json:"invite_code"`
	AnniversaryDate *time.Time `gorm:"type:date" json:"anniversary_date"`
	Status          string     `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending | active | ended
	CreatedAt       time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relasi ke User
	UserOne User  `gorm:"foreignKey:UserOneID" json:"user_one"`
	UserTwo *User `gorm:"foreignKey:UserTwoID" json:"user_two"`
}