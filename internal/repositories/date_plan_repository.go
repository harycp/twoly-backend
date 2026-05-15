package repositories

import (
	"github.com/harycp/twoly-backend/internal/models"
	"gorm.io/gorm"
)

type DatePlanRepository interface {
	CreateDatePlan(datePlan *models.DatePlan) error
	FindAllByCoupleID(coupleID string, status string) ([]models.DatePlan, error)
	FindByID(id string, coupleID string) (*models.DatePlan, error)
	UpdateDatePlan(datePlan *models.DatePlan) error
	DeleteDatePlan(datePlan *models.DatePlan) error
	FindChecklistByID(checklistID string, planID string) (*models.DatePlanChecklist, error)
	UpdateChecklist(checklist *models.DatePlanChecklist) error
}

type datePlanRepository struct {
	db *gorm.DB
}

func NewDatePlanRepository(db *gorm.DB) DatePlanRepository {
	return &datePlanRepository{db}
}

func (r *datePlanRepository) CreateDatePlan(datePlan *models.DatePlan) error {
	return r.db.Create(datePlan).Error
}

func (r *datePlanRepository) FindAllByCoupleID(coupleID string, status string) ([]models.DatePlan, error) {
	var datePlans []models.DatePlan
	query := r.db.Preload("Checklists").Where("couple_id = ?", coupleID)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	// Order by closest plan date first
	err := query.Order("plan_date ASC").Find(&datePlans).Error
	return datePlans, err
}

func (r *datePlanRepository) FindByID(id string, coupleID string) (*models.DatePlan, error) {
	var datePlan models.DatePlan
	err := r.db.Preload("Checklists").Where("id = ? AND couple_id = ?", id, coupleID).First(&datePlan).Error
	if err != nil {
		return nil, err
	}
	return &datePlan, nil
}

func (r *datePlanRepository) UpdateDatePlan(datePlan *models.DatePlan) error {
	return r.db.Save(datePlan).Error
}

func (r *datePlanRepository) DeleteDatePlan(datePlan *models.DatePlan) error {
	return r.db.Delete(datePlan).Error
}

func (r *datePlanRepository) FindChecklistByID(checklistID string, planID string) (*models.DatePlanChecklist, error) {
	var checklist models.DatePlanChecklist
	err := r.db.Where("id = ? AND date_plan_id = ?", checklistID, planID).First(&checklist).Error
	if err != nil {
		return nil, err
	}
	return &checklist, nil
}

func (r *datePlanRepository) UpdateChecklist(checklist *models.DatePlanChecklist) error {
	return r.db.Save(checklist).Error
}