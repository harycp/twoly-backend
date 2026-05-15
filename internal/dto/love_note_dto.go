package dto

import "time"

type CreateLoveNoteRequest struct {
	Message  string `json:"message" binding:"required"`
	UnlockAt string `json:"unlock_at" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"` // ISO 8601 format
}

type LoveNoteResponse struct {
	ID         string     `json:"id"`
	CoupleID   string     `json:"couple_id"`
	SenderID   string     `json:"sender_id"`
	ReceiverID string     `json:"receiver_id"`
	Message    string     `json:"message"` // Teks ini akan di-hidden (disensor) oleh Service jika belum waktunya dibuka
	UnlockAt   *time.Time `json:"unlock_at"`
	IsOpened   bool       `json:"is_opened"`
	OpenedAt   *time.Time `json:"opened_at"`
	CreatedAt  time.Time  `json:"created_at"`
}