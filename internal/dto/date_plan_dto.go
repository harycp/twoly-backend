package dto

import "time"

type CreateDatePlanRequest struct {
	Title        string   `json:"title" binding:"required"`
	Notes        string   `json:"notes"`
	PlanDate     string   `json:"plan_date" binding:"required,datetime=2006-01-02T15:04:05Z07:00"` // Membutuhkan format ISO 8601 dengan timezone
	LocationName string   `json:"location_name"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	Checklists   []string `json:"checklists"` // Mengirim daftar checklist saat membuat date plan
}

type UpdateDatePlanRequest struct {
	Title        string   `json:"title"`
	Notes        string   `json:"notes"`
	PlanDate     string   `json:"plan_date" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	LocationName string   `json:"location_name"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
}

type UpdateDatePlanStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=planned ongoing completed cancelled"`
}

type UpdateChecklistItemRequest struct {
	IsChecked bool `json:"is_checked"`
}

type AddChecklistItemRequest struct {
	Item string `json:"item" binding:"required"`
}

type ChecklistItemResponse struct {
	ID        string `json:"id"`
	Item      string `json:"item"`
	IsChecked bool   `json:"is_checked"`
}

type DatePlanResponse struct {
	ID           string                  `json:"id"`
	CoupleID     string                  `json:"couple_id"`
	CreatedBy    string                  `json:"created_by"`
	Title        string                  `json:"title"`
	Notes        string                  `json:"notes"`
	PlanDate     time.Time               `json:"plan_date"`
	LocationName string                  `json:"location_name"`
	Latitude     *float64                `json:"latitude"`
	Longitude    *float64                `json:"longitude"`
	Status       string                  `json:"status"`
	Checklists   []ChecklistItemResponse `json:"checklists"`
	CreatedAt    time.Time               `json:"created_at"`
}