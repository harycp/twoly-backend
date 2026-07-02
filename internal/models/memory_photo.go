package models

import (
	"time"

	"github.com/google/uuid"
)

type MemoryPhoto struct {
	ID                 uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	MemoryID           uuid.UUID `gorm:"type:uuid;not null" json:"memory_id"`
	UploadedBy         uuid.UUID `gorm:"type:uuid;not null" json:"uploaded_by"`
	MediaType          string    `gorm:"type:varchar(20);not null;default:'image'" json:"media_type"`
	PhotoURL           string    `gorm:"type:text;not null" json:"photo_url"`
	CloudinaryPublicID string    `gorm:"type:text" json:"cloudinary_public_id"`
	Caption            string    `gorm:"type:text" json:"caption"`
	CreatedAt          time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	Memory   Memory `gorm:"foreignKey:MemoryID" json:"-"`
	Uploader User   `gorm:"foreignKey:UploadedBy" json:"-"`
}