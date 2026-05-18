package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name           string         `gorm:"type:varchar(100);not null" json:"name"`
	Username       string         `gorm:"type:varchar(50);unique;not null" json:"username"`
	Email          string         `gorm:"type:varchar(150);unique;not null" json:"email"`
	PasswordHash   string         `gorm:"type:text;not null" json:"-"`
	AvatarURL      *string        `gorm:"type:text" json:"avatar_url"`
	AvatarPublicID *string        `gorm:"column:avatar_cloudinary_public_id;type:text" json:"-"`
	CreatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
