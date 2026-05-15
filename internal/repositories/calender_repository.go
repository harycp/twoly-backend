package repositories

import (
	"time"

	"github.com/harycp/twoly-backend/internal/models"
	"gorm.io/gorm"
)

type CalendarRepository interface {
	CreateCustomEvent(event *models.CalendarEvent) error
	FindCustomEventsByDateRange(coupleID string, start time.Time, end time.Time) ([]models.CalendarEvent, error)
	FindByID(id string, coupleID string) (*models.CalendarEvent, error)
	UpdateCustomEvent(event *models.CalendarEvent) error
	DeleteCustomEvent(event *models.CalendarEvent) error
}

type calendarRepository struct {
	db *gorm.DB
}

func NewCalendarRepository(db *gorm.DB) CalendarRepository {
	return &calendarRepository{db}
}

func (r *calendarRepository) CreateCustomEvent(event *models.CalendarEvent) error {
	return r.db.Create(event).Error
}

func (r *calendarRepository) FindCustomEventsByDateRange(coupleID string, start time.Time, end time.Time) ([]models.CalendarEvent, error) {
	var events []models.CalendarEvent
	err := r.db.Where("couple_id = ? AND event_date >= ? AND event_date <= ?", coupleID, start, end).
		Order("event_date ASC").Find(&events).Error
	return events, err
}

func (r *calendarRepository) FindByID(id string, coupleID string) (*models.CalendarEvent, error) {
	var event models.CalendarEvent
	err := r.db.Where("id = ? AND couple_id = ?", id, coupleID).First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *calendarRepository) UpdateCustomEvent(event *models.CalendarEvent) error {
	return r.db.Save(event).Error
}

func (r *calendarRepository) DeleteCustomEvent(event *models.CalendarEvent) error {
	return r.db.Delete(event).Error
}