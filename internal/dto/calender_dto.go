package dto

import "time"

type CreateCustomEventRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	EventDate   string `json:"event_date" binding:"required,datetime=2006-01-02T15:04:05Z07:00"` // ISO 8601
	EventType   string `json:"event_type" binding:"required,oneof=custom"`
}

type UpdateCustomEventRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	EventDate   string `json:"event_date" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"` // ISO 8601
}

type GetCalendarEventsRequest struct {
	StartDate string `form:"start_date" binding:"required,datetime=2006-01-02"` 
	EndDate   string `form:"end_date" binding:"required,datetime=2006-01-02"`  
}

type CalendarEventResponse struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description,omitempty"`
	EventDate    time.Time `json:"event_date"`
	EventType    string    `json:"event_type"`
	ReferenceID  string    `json:"reference_id,omitempty"`
	
	LocationName string    `json:"location_name,omitempty"`
	Status       string    `json:"status,omitempty"`       
}