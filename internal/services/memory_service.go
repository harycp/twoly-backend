package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/harycp/twoly-backend/internal/dto"
	"github.com/harycp/twoly-backend/internal/models"
	"github.com/harycp/twoly-backend/internal/repositories"
	"github.com/lib/pq"
)

type MemoryService interface {
	CreateMemory(userID string, req dto.CreateMemoryRequest) (dto.MemoryResponse, error)
	GetAllMemories(userID string, month string) ([]dto.MemoryResponse, error)
	GetMemoryDetail(userID string, memoryID string) (dto.MemoryResponse, error)
	UpdateMemory(userID string, memoryID string, req dto.UpdateMemoryRequest) (dto.MemoryResponse, error)
	DeleteMemory(userID string, memoryID string) error
}

type memoryService struct {
	memoryRepo repositories.MemoryRepository
	coupleRepo repositories.CoupleRepository
}

func NewMemoryService(memoryRepo repositories.MemoryRepository, coupleRepo repositories.CoupleRepository) MemoryService {
	return &memoryService{memoryRepo, coupleRepo}
}

func (s *memoryService) getActiveCoupleID(userID string) (uuid.UUID, error) {
	couple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil || couple == nil {
		return uuid.Nil, errors.New("you are not part of any couple")
	}
	if couple.Status != "active" {
		return uuid.Nil, errors.New("your couple status is not active yet")
	}
	return couple.ID, nil
}

func (s *memoryService) CreateMemory(userID string, req dto.CreateMemoryRequest) (dto.MemoryResponse, error) {
	coupleID, err := s.getActiveCoupleID(userID)
	if err != nil {
		return dto.MemoryResponse{}, err
	}

	userUUID, _ := uuid.Parse(userID)
	memoryDate, _ := time.Parse("2006-01-02", req.MemoryDate)

	memory := models.Memory{
		CoupleID:     coupleID,
		CreatedBy:    userUUID,
		Title:        req.Title,
		Description:  req.Description,
		MemoryDate:   memoryDate,
		LocationName: req.LocationName,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Mood:         req.Mood,
		Tags:         pq.StringArray(req.Tags),
	}

	err = s.memoryRepo.CreateMemory(&memory)
	if err != nil {
		return dto.MemoryResponse{}, err
	}

	return s.mapToResponse(memory), nil
}

func (s *memoryService) GetAllMemories(userID string, month string) ([]dto.MemoryResponse, error) {
	coupleID, err := s.getActiveCoupleID(userID)
	if err != nil {
		return nil, err
	}

	memories, err := s.memoryRepo.FindAllByCoupleID(coupleID.String(), month)
	if err != nil {
		return nil, err
	}

	var response []dto.MemoryResponse
	for _, m := range memories {
		response = append(response, s.mapToResponse(m))
	}

	if response == nil {
		response = []dto.MemoryResponse{}
	}

	return response, nil
}

func (s *memoryService) GetMemoryDetail(userID string, memoryID string) (dto.MemoryResponse, error) {
	coupleID, err := s.getActiveCoupleID(userID)
	if err != nil {
		return dto.MemoryResponse{}, err
	}

	memory, err := s.memoryRepo.FindByID(memoryID, coupleID.String())
	if err != nil || memory == nil {
		return dto.MemoryResponse{}, errors.New("memory not found")
	}

	return s.mapToResponse(*memory), nil
}

func (s *memoryService) UpdateMemory(userID string, memoryID string, req dto.UpdateMemoryRequest) (dto.MemoryResponse, error) {
	coupleID, err := s.getActiveCoupleID(userID)
	if err != nil {
		return dto.MemoryResponse{}, err
	}

	memory, err := s.memoryRepo.FindByID(memoryID, coupleID.String())
	if err != nil || memory == nil {
		return dto.MemoryResponse{}, errors.New("memory not found")
	}

	if req.Title != "" {
		memory.Title = req.Title
	}
	if req.Description != "" {
		memory.Description = req.Description
	}
	if req.MemoryDate != "" {
		memoryDate, _ := time.Parse("2006-01-02", req.MemoryDate)
		memory.MemoryDate = memoryDate
	}
	if req.LocationName != "" {
		memory.LocationName = req.LocationName
	}
	if req.Latitude != nil {
		memory.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		memory.Longitude = req.Longitude
	}
	if req.Mood != "" {
		memory.Mood = req.Mood
	}
	if len(req.Tags) > 0 {
		memory.Tags = pq.StringArray(req.Tags)
	}

	err = s.memoryRepo.UpdateMemory(memory)
	if err != nil {
		return dto.MemoryResponse{}, err
	}

	return s.mapToResponse(*memory), nil
}

func (s *memoryService) DeleteMemory(userID string, memoryID string) error {
	coupleID, err := s.getActiveCoupleID(userID)
	if err != nil {
		return err
	}

	memory, err := s.memoryRepo.FindByID(memoryID, coupleID.String())
	if err != nil || memory == nil {
		return errors.New("memory not found")
	}

	return s.memoryRepo.DeleteMemory(memory)
}

func (s *memoryService) mapToResponse(m models.Memory) dto.MemoryResponse {
	return dto.MemoryResponse{
		ID:           m.ID,
		CoupleID:     m.CoupleID,
		CreatedBy:    m.CreatedBy,
		Title:        m.Title,
		Description:  m.Description,
		MemoryDate:   m.MemoryDate,
		LocationName: m.LocationName,
		Latitude:     m.Latitude,
		Longitude:    m.Longitude,
		Mood:         m.Mood,
		Tags:         m.Tags,
		CreatedAt:    m.CreatedAt,
	}
}