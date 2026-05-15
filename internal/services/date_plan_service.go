package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/harycp/twoly-backend/internal/dto"
	"github.com/harycp/twoly-backend/internal/models"
	"github.com/harycp/twoly-backend/internal/repositories"
)

type DatePlanService interface {
	CreateDatePlan(userID string, req dto.CreateDatePlanRequest) (dto.DatePlanResponse, error)
	GetAllDatePlans(userID string, status string) ([]dto.DatePlanResponse, error)
	GetDatePlanDetail(userID string, planID string) (dto.DatePlanResponse, error)
	UpdateDatePlan(userID string, planID string, req dto.UpdateDatePlanRequest) (dto.DatePlanResponse, error)
	DeleteDatePlan(userID string, planID string) error
	UpdateDatePlanStatus(userID string, planID string, req dto.UpdateDatePlanStatusRequest) (dto.DatePlanResponse, error)
	ConvertToMemory(userID string, planID string) (dto.MemoryResponse, error)
	UpdateChecklistItem(userID string, planID string, checklistID string, req dto.UpdateChecklistItemRequest) (dto.ChecklistItemResponse, error)
}

type datePlanService struct {
	datePlanRepo repositories.DatePlanRepository
	coupleRepo   repositories.CoupleRepository
	memoryRepo   repositories.MemoryRepository
}

func NewDatePlanService(
	datePlanRepo repositories.DatePlanRepository,
	coupleRepo   repositories.CoupleRepository,
	memoryRepo   repositories.MemoryRepository,
) DatePlanService {
	return &datePlanService{datePlanRepo, coupleRepo, memoryRepo}
}

func (s *datePlanService) getActiveCoupleID(userID string) (uuid.UUID, error) {
	couple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil || couple == nil || couple.Status != "active" {
		return uuid.Nil, errors.New("you do not have an active couple")
	}
	return couple.ID, nil
}

func (s *datePlanService) CreateDatePlan(userID string, req dto.CreateDatePlanRequest) (dto.DatePlanResponse, error) {
	coupleID, err := s.getActiveCoupleID(userID)
	if err != nil {
		return dto.DatePlanResponse{}, err
	}

	userUUID, _ := uuid.Parse(userID)
	
	// Parse date based on ISO 8601
	planDate, err := time.Parse(time.RFC3339, req.PlanDate)
	if err != nil {
		return dto.DatePlanResponse{}, errors.New("invalid date format, please use ISO 8601 (e.g. 2026-05-20T19:00:00+07:00)")
	}

	var checklists []models.DatePlanChecklist
	for _, item := range req.Checklists {
		if item != "" {
			checklists = append(checklists, models.DatePlanChecklist{
				Item: item,
			})
		}
	}

	datePlan := models.DatePlan{
		CoupleID:     coupleID,
		CreatedBy:    userUUID,
		Title:        req.Title,
		Notes:        req.Notes,
		PlanDate:     planDate,
		LocationName: req.LocationName,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Status:       "planned", // default status
		Checklists:   checklists,
	}

	if err := s.datePlanRepo.CreateDatePlan(&datePlan); err != nil {
		return dto.DatePlanResponse{}, errors.New("failed to create date plan")
	}

	return s.mapToResponse(datePlan), nil
}

func (s *datePlanService) GetAllDatePlans(userID string, status string) ([]dto.DatePlanResponse, error) {
	coupleID, err := s.getActiveCoupleID(userID)
	if err != nil {
		return nil, err
	}

	plans, err := s.datePlanRepo.FindAllByCoupleID(coupleID.String(), status)
	if err != nil {
		return nil, err
	}

	var response []dto.DatePlanResponse
	for _, p := range plans {
		response = append(response, s.mapToResponse(p))
	}

	if response == nil {
		response = []dto.DatePlanResponse{}
	}

	return response, nil
}

func (s *datePlanService) GetDatePlanDetail(userID string, planID string) (dto.DatePlanResponse, error) {
	coupleID, err := s.getActiveCoupleID(userID)
	if err != nil {
		return dto.DatePlanResponse{}, err
	}

	plan, err := s.datePlanRepo.FindByID(planID, coupleID.String())
	if err != nil || plan == nil {
		return dto.DatePlanResponse{}, errors.New("date plan not found")
	}

	return s.mapToResponse(*plan), nil
}

func (s *datePlanService) UpdateDatePlan(userID string, planID string, req dto.UpdateDatePlanRequest) (dto.DatePlanResponse, error) {
	coupleID, err := s.getActiveCoupleID(userID)
	if err != nil {
		return dto.DatePlanResponse{}, err
	}

	plan, err := s.datePlanRepo.FindByID(planID, coupleID.String())
	if err != nil || plan == nil {
		return dto.DatePlanResponse{}, errors.New("date plan not found")
	}

	if req.Title != "" {
		plan.Title = req.Title
	}
	if req.Notes != "" {
		plan.Notes = req.Notes
	}
	if req.PlanDate != "" {
		planDate, err := time.Parse(time.RFC3339, req.PlanDate)
		if err == nil {
			plan.PlanDate = planDate
		}
	}
	if req.LocationName != "" {
		plan.LocationName = req.LocationName
	}
	if req.Latitude != nil {
		plan.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		plan.Longitude = req.Longitude
	}

	if err := s.datePlanRepo.UpdateDatePlan(plan); err != nil {
		return dto.DatePlanResponse{}, errors.New("failed to update date plan")
	}

	return s.mapToResponse(*plan), nil
}

func (s *datePlanService) DeleteDatePlan(userID string, planID string) error {
	coupleID, err := s.getActiveCoupleID(userID)
	if err != nil {
		return err
	}

	plan, err := s.datePlanRepo.FindByID(planID, coupleID.String())
	if err != nil || plan == nil {
		return errors.New("date plan not found")
	}

	return s.datePlanRepo.DeleteDatePlan(plan)
}

func (s *datePlanService) UpdateDatePlanStatus(userID string, planID string, req dto.UpdateDatePlanStatusRequest) (dto.DatePlanResponse, error) {
	coupleID, err := s.getActiveCoupleID(userID)
	if err != nil {
		return dto.DatePlanResponse{}, err
	}

	plan, err := s.datePlanRepo.FindByID(planID, coupleID.String())
	if err != nil || plan == nil {
		return dto.DatePlanResponse{}, errors.New("date plan not found")
	}

	plan.Status = req.Status
	if err := s.datePlanRepo.UpdateDatePlan(plan); err != nil {
		return dto.DatePlanResponse{}, errors.New("failed to update date plan status")
	}

	return s.mapToResponse(*plan), nil
}

func (s *datePlanService) ConvertToMemory(userID string, planID string) (dto.MemoryResponse, error) {
	coupleID, err := s.getActiveCoupleID(userID)
	if err != nil {
		return dto.MemoryResponse{}, err
	}

	plan, err := s.datePlanRepo.FindByID(planID, coupleID.String())
	if err != nil || plan == nil {
		return dto.MemoryResponse{}, errors.New("date plan not found")
	}

	if plan.Status != "completed" {
		plan.Status = "completed"
		_ = s.datePlanRepo.UpdateDatePlan(plan)
	}

	userUUID, _ := uuid.Parse(userID)
	
	// Create a new Memory record using data from DatePlan
	memory := models.Memory{
		CoupleID:     coupleID,
		CreatedBy:    userUUID,
		Title:        plan.Title,
		Description:  plan.Notes,
		MemoryDate:   plan.PlanDate, 
		LocationName: plan.LocationName,
		Latitude:     plan.Latitude,
		Longitude:    plan.Longitude,
	}

	if err := s.memoryRepo.CreateMemory(&memory); err != nil {
		return dto.MemoryResponse{}, errors.New("failed to create memory from date plan")
	}

	return dto.MemoryResponse{
		ID:           memory.ID,
		CoupleID:     memory.CoupleID,
		CreatedBy:    memory.CreatedBy,
		Title:        memory.Title,
		Description:  memory.Description,
		MemoryDate:   memory.MemoryDate,
		LocationName: memory.LocationName,
		Latitude:     memory.Latitude,
		Longitude:    memory.Longitude,
		Mood:         memory.Mood,
		CreatedAt:    memory.CreatedAt,
	}, nil
}

func (s *datePlanService) UpdateChecklistItem(userID string, planID string, checklistID string, req dto.UpdateChecklistItemRequest) (dto.ChecklistItemResponse, error) {
	coupleID, err := s.getActiveCoupleID(userID)
	if err != nil {
		return dto.ChecklistItemResponse{}, err
	}

	plan, err := s.datePlanRepo.FindByID(planID, coupleID.String())
	if err != nil || plan == nil {
		return dto.ChecklistItemResponse{}, errors.New("date plan not found")
	}

	checklist, err := s.datePlanRepo.FindChecklistByID(checklistID, planID)
	if err != nil || checklist == nil {
		return dto.ChecklistItemResponse{}, errors.New("checklist item not found")
	}

	checklist.IsChecked = req.IsChecked
	if err := s.datePlanRepo.UpdateChecklist(checklist); err != nil {
		return dto.ChecklistItemResponse{}, errors.New("failed to update checklist item")
	}

	return dto.ChecklistItemResponse{
		ID:        checklist.ID.String(),
		Item:      checklist.Item,
		IsChecked: checklist.IsChecked,
	}, nil
}

func (s *datePlanService) mapToResponse(p models.DatePlan) dto.DatePlanResponse {
	var checklists []dto.ChecklistItemResponse
	for _, c := range p.Checklists {
		checklists = append(checklists, dto.ChecklistItemResponse{
			ID:        c.ID.String(),
			Item:      c.Item,
			IsChecked: c.IsChecked,
		})
	}

	if checklists == nil {
		checklists = []dto.ChecklistItemResponse{}
	}

	return dto.DatePlanResponse{
		ID:           p.ID.String(),
		CoupleID:     p.CoupleID.String(),
		CreatedBy:    p.CreatedBy.String(),
		Title:        p.Title,
		Notes:        p.Notes,
		PlanDate:     p.PlanDate,
		LocationName: p.LocationName,
		Latitude:     p.Latitude,
		Longitude:    p.Longitude,
		Status:       p.Status,
		Checklists:   checklists,
		CreatedAt:    p.CreatedAt,
	}
}