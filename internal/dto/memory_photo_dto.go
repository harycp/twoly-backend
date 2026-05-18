package dto

import "time"

type MemoryPhotoResponse struct {
	ID                 string    `json:"id"`
	MemoryID           string    `json:"memory_id"`
	UploadedBy         string    `json:"uploaded_by"`
	PhotoURL           string    `json:"photo_url"`
	CloudinaryPublicID string    `json:"cloudinary_public_id"`
	Caption            string    `json:"caption"`
	CreatedAt          time.Time `json:"created_at"`
}

type GalleryMemoryDetail struct {
	ID           string    `json:"id"`
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

type GalleryPhotoResponse struct {
	ID                 string              `json:"id"`
	MemoryID           string              `json:"memory_id"`
	UploadedBy         string              `json:"uploaded_by"`
	PhotoURL           string              `json:"photo_url"`
	CloudinaryPublicID string              `json:"cloudinary_public_id"`
	Caption            string              `json:"caption"`
	CreatedAt          time.Time           `json:"created_at"`
	Memory             GalleryMemoryDetail `json:"memory"`
}