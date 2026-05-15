package services

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/harycp/twoly-backend/internal/dto"
	"github.com/harycp/twoly-backend/internal/models"
	"github.com/harycp/twoly-backend/internal/repositories"
)

type CalendarService interface {
	CreateCustomEvent(userID string, req dto.CreateCustomEventRequest) (dto.CalendarEventResponse, error)
	GetCalendarEvents(userID string, start string, end string) ([]dto.CalendarEventResponse, error)
	UpdateCustomEvent(userID string, eventID string, req dto.UpdateCustomEventRequest) (dto.CalendarEventResponse, error)
	DeleteCustomEvent(userID string, eventID string) error
}

type calendarService struct {
	calendarRepo repositories.CalendarRepository
	memoryRepo   repositories.MemoryRepository
	datePlanRepo repositories.DatePlanRepository
	coupleRepo   repositories.CoupleRepository
}

func NewCalendarService(
	calendarRepo repositories.CalendarRepository,
	memoryRepo repositories.MemoryRepository,
	datePlanRepo repositories.DatePlanRepository,
	coupleRepo repositories.CoupleRepository,
) CalendarService {
	return &calendarService{calendarRepo, memoryRepo, datePlanRepo, coupleRepo}
}

func (s *calendarService) getActiveCouple(userID string) (*models.Couple, error) {
	couple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil || couple == nil || couple.Status != "active" {
		return nil, errors.New("you do not have an active couple")
	}
	return couple, nil
}

func (s *calendarService) CreateCustomEvent(userID string, req dto.CreateCustomEventRequest) (dto.CalendarEventResponse, error) {
	couple, err := s.getActiveCouple(userID)
	if err != nil {
		return dto.CalendarEventResponse{}, err
	}

	userUUID, _ := uuid.Parse(userID)
	eventDate, _ := time.Parse(time.RFC3339, req.EventDate)

	event := models.CalendarEvent{
		CoupleID:    couple.ID,
		CreatedBy:   userUUID,
		Title:       req.Title,
		Description: req.Description,
		EventDate:   eventDate,
		EventType:   req.EventType,
	}

	if err := s.calendarRepo.CreateCustomEvent(&event); err != nil {
		return dto.CalendarEventResponse{}, errors.New("failed to create custom event")
	}

	return s.mapCustomToResponse(event), nil
}

func (s *calendarService) GetCalendarEvents(userID string, startStr string, endStr string) ([]dto.CalendarEventResponse, error) {
	couple, err := s.getActiveCouple(userID)
	if err != nil {
		return nil, err
	}

	// Parse boundary dates
	startDate, errStart := time.Parse("2006-01-02", startStr)
	endDate, errEnd := time.Parse("2006-01-02", endStr)
	if errStart != nil || errEnd != nil {
		return nil, errors.New("invalid date format, please use YYYY-MM-DD for start and end queries")
	}

	// Adjust endDate to include the entire day until 23:59:59
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	var allEvents []dto.CalendarEventResponse

	// 1. Fetch Custom Events
	customEvents, _ := s.calendarRepo.FindCustomEventsByDateRange(couple.ID.String(), startDate, endDate)
	for _, ce := range customEvents {
		allEvents = append(allEvents, s.mapCustomToResponse(ce))
	}

	// 2. Fetch Memories
	memories, _ := s.memoryRepo.FindAllByDateRange(couple.ID.String(), startDate, endDate)
	for _, m := range memories {
		allEvents = append(allEvents, dto.CalendarEventResponse{
			ID:           m.ID.String(),
			Title:        m.Title,
			Description:  m.Description,
			EventDate:    m.MemoryDate,
			EventType:    "memory",
			ReferenceID:  m.ID.String(),
			LocationName: m.LocationName,
		})
	}

	// 3. Fetch Date Plans
	datePlans, _ := s.datePlanRepo.FindAllByDateRange(couple.ID.String(), startDate, endDate)
	for _, dp := range datePlans {
		allEvents = append(allEvents, dto.CalendarEventResponse{
			ID:           dp.ID.String(),
			Title:        dp.Title,
			Description:  dp.Notes,
			EventDate:    dp.PlanDate,
			EventType:    "date_plan",
			ReferenceID:  dp.ID.String(),
			LocationName: dp.LocationName,
			Status:       dp.Status,
		})
	}

	if couple.AnniversaryDate != nil {
		annivMonth := couple.AnniversaryDate.Month()
		annivDay := couple.AnniversaryDate.Day()
		startYear := startDate.Year()
		endYear := endDate.Year()

		for year := startYear; year <= endYear; year++ {
			annivThisYear := time.Date(year, annivMonth, annivDay, 0, 0, 0, 0, startDate.Location())
			
			if (annivThisYear.After(startDate) || annivThisYear.Equal(startDate)) && 
			(annivThisYear.Before(endDate) || annivThisYear.Equal(endDate)) {
				
				allEvents = append(allEvents, dto.CalendarEventResponse{
					ID:        fmt.Sprintf("%s-anniv-%d", couple.ID.String(), year),
					Title:     "Anniversary",
					EventDate: annivThisYear,
					EventType: "anniversary",
				})
			}
		}
	}

	// 5. Sort all combined events by EventDate Ascending
	sort.Slice(allEvents, func(i, j int) bool {
		return allEvents[i].EventDate.Before(allEvents[j].EventDate)
	})

	if allEvents == nil {
		allEvents = []dto.CalendarEventResponse{}
	}

	return allEvents, nil
}

func (s *calendarService) UpdateCustomEvent(userID string, eventID string, req dto.UpdateCustomEventRequest) (dto.CalendarEventResponse, error) {
	couple, err := s.getActiveCouple(userID)
	if err != nil {
		return dto.CalendarEventResponse{}, err
	}

	event, err := s.calendarRepo.FindByID(eventID, couple.ID.String())
	if err != nil || event == nil {
		return dto.CalendarEventResponse{}, errors.New("custom event not found")
	}

	if req.Title != "" {
		event.Title = req.Title
	}
	if req.Description != "" {
		event.Description = req.Description
	}
	if req.EventDate != "" {
		eventDate, err := time.Parse(time.RFC3339, req.EventDate)
		if err == nil {
			event.EventDate = eventDate
		}
	}

	if err := s.calendarRepo.UpdateCustomEvent(event); err != nil {
		return dto.CalendarEventResponse{}, errors.New("failed to update event")
	}

	return s.mapCustomToResponse(*event), nil
}

func (s *calendarService) DeleteCustomEvent(userID string, eventID string) error {
	couple, err := s.getActiveCouple(userID)
	if err != nil {
		return err
	}

	event, err := s.calendarRepo.FindByID(eventID, couple.ID.String())
	if err != nil || event == nil {
		return errors.New("custom event not found")
	}

	return s.calendarRepo.DeleteCustomEvent(event)
}

func (s *calendarService) mapCustomToResponse(e models.CalendarEvent) dto.CalendarEventResponse {
	refID := ""
	if e.ReferenceID != nil {
		refID = e.ReferenceID.String()
	}
	return dto.CalendarEventResponse{
		ID:          e.ID.String(),
		Title:       e.Title,
		Description: e.Description,
		EventDate:   e.EventDate,
		EventType:   e.EventType,
		ReferenceID: refID,
	}
}