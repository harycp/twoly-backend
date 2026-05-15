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