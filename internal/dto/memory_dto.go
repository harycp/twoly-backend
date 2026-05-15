package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateMemoryRequest struct {
	Title        string   `json:"title" binding:"required"`
	Description  string   `json:"description"`
	MemoryDate   string   `json:"memory_date" binding:"required,datetime=2006-01-02"`
	LocationName string   `json:"location_name"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	Mood         string   `json:"mood"`
	Tags         []string `json:"tags"`
}

type UpdateMemoryRequest struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	MemoryDate   string   `json:"memory_date" binding:"omitempty,datetime=2006-01-02"`
	LocationName string   `json:"location_name"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	Mood         string   `json:"mood"`
	Tags         []string `json:"tags"`
}

type MemoryResponse struct {
	ID           uuid.UUID `json:"id"`
	CoupleID     uuid.UUID `json:"couple_id"`
	CreatedBy    uuid.UUID `json:"created_by"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	MemoryDate   time.Time `json:"memory_date"`
	LocationName string    `json:"location_name"`
	Latitude     *float64  `json:"latitude"`
	Longitude    *float64  `json:"longitude"`
	Mood         string    `json:"mood"`
	Tags         []string  `json:"tags"`
	CreatedAt    time.Time `json:"created_at"`
}