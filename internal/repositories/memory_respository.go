package repositories

import (
	"github.com/harycp/twoly-backend/internal/models"
	"gorm.io/gorm"
)

type MemoryRepository interface {
	CreateMemory(memory *models.Memory) error
	FindAllByCoupleID(coupleID string) ([]models.Memory, error)
	FindByID(id string, coupleID string) (*models.Memory, error)
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

func (r *memoryRepository) FindAllByCoupleID(coupleID string) ([]models.Memory, error) {
	var memories []models.Memory
	err := r.db.Where("couple_id = ?", coupleID).Order("memory_date DESC").Find(&memories).Error
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

func (r *memoryRepository) UpdateMemory(memory *models.Memory) error {
	return r.db.Save(memory).Error
}

func (r *memoryRepository) DeleteMemory(memory *models.Memory) error {
	return r.db.Delete(memory).Error
}