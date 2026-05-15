package models

import (
	"time"

	"github.com/google/uuid"
)

type LoveNote struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CoupleID   uuid.UUID  `gorm:"type:uuid;not null" json:"couple_id"`
	SenderID   uuid.UUID  `gorm:"type:uuid;not null" json:"sender_id"`
	ReceiverID uuid.UUID  `gorm:"type:uuid;not null" json:"receiver_id"`
	Message    string     `gorm:"type:text;not null" json:"message"`
	UnlockAt   *time.Time `gorm:"type:timestamp" json:"unlock_at"`
	IsOpened   bool       `gorm:"default:false" json:"is_opened"`
	OpenedAt   *time.Time `gorm:"type:timestamp" json:"opened_at"`
	CreatedAt  time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	Couple   Couple `gorm:"foreignKey:CoupleID" json:"-"`
	Sender   User   `gorm:"foreignKey:SenderID" json:"-"`
	Receiver User   `gorm:"foreignKey:ReceiverID" json:"-"`
}