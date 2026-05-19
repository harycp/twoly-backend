package repositories

import (
	"time"

	"github.com/harycp/twoly-backend/internal/models"
	"gorm.io/gorm"
)

type MemoryRepository interface {
	CreateMemory(memory *models.Memory) error
	FindAllByCoupleID(coupleID string, month string) ([]models.Memory, error)
	FindAllByDateRange(coupleID string, start time.Time, end time.Time) ([]models.Memory, error)
	FindByID(id string, coupleID string) (*models.Memory, error)
	FindByConvertedFromDatePlanID(datePlanID string) (*models.Memory, error)
	UpdateMemory(memory *models.Memory) error
	DeleteMemory(memory *models.Memory) error
}

type memoryRepository struct {
	db *gorm.DB
}

func NewMemoryRepository(db *gorm.DB) MemoryRepository {
	return &memoryRepository{db}
}

func (r *memoryRepository) CreateMemory(memory *models.Memory) error {
	return r.db.Create(memory).Error
}

func (r *memoryRepository) FindAllByCoupleID(coupleID string, month string) ([]models.Memory, error) {
	var memories []models.Memory
	query := r.db.Where("couple_id = ?", coupleID)
	
	if month != "" {
		query = query.Where("TO_CHAR(memory_date, 'YYYY-MM') = ?", month)
	}
	
	err := query.Order("memory_date ASC").Find(&memories).Error
	return memories, err
}

func (r *memoryRepository) FindAllByDateRange(coupleID string, start time.Time, end time.Time) ([]models.Memory, error) {
	var memories []models.Memory
	err := r.db.Where("couple_id = ? AND memory_date >= ? AND memory_date <= ?", coupleID, start, end).
		Order("memory_date ASC").Find(&memories).Error
	return memories, err
}

func (r *memoryRepository) FindByID(id string, coupleID string) (*models.Memory, error) {
	var memory models.Memory
	err := r.db.Where("id = ? AND couple_id = ?", id, coupleID).First(&memory).Error
	if err != nil {
		return nil, err
	}
	return &memory, nil
}

func (r *memoryRepository) FindByConvertedFromDatePlanID(datePlanID string) (*models.Memory, error) {
	var memory models.Memory
	err := r.db.Where("converted_from_date_plan_id = ?", datePlanID).First(&memory).Error
	if err != nil {
		return nil, err
	}
	return &memory, nil
}

func (r *memoryRepository) UpdateMemory(memory *models.Memory) error {
	return r.db.Save(memory).Error
}

func (r *memoryRepository) DeleteMemory(memory *models.Memory) error {
	return r.db.Delete(memory).Error
}