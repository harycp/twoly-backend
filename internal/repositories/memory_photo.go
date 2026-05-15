package repositories

import (
	"github.com/harycp/twoly-backend/internal/models"
	"gorm.io/gorm"
)

type MemoryPhotoRepository interface {
	CreatePhoto(photo *models.MemoryPhoto) error
	FindByMemoryID(memoryID string) ([]models.MemoryPhoto, error)
	FindByID(id string) (*models.MemoryPhoto, error)
	DeletePhoto(photo *models.MemoryPhoto) error
}

type memoryPhotoRepository struct {
	db *gorm.DB
}

func NewMemoryPhotoRepository(db *gorm.DB) MemoryPhotoRepository {
	return &memoryPhotoRepository{db}
}

func (r *memoryPhotoRepository) CreatePhoto(photo *models.MemoryPhoto) error {
	return r.db.Create(photo).Error
}

func (r *memoryPhotoRepository) FindByMemoryID(memoryID string) ([]models.MemoryPhoto, error) {
	var photos []models.MemoryPhoto
	err := r.db.Where("memory_id = ?", memoryID).Order("created_at DESC").Find(&photos).Error
	return photos, err
}

func (r *memoryPhotoRepository) FindByID(id string) (*models.MemoryPhoto, error) {
	var photo models.MemoryPhoto
	err := r.db.Where("id = ?", id).First(&photo).Error
	if err != nil {
		return nil, err
	}
	return &photo, nil
}

func (r *memoryPhotoRepository) DeletePhoto(photo *models.MemoryPhoto) error {
	return r.db.Delete(photo).Error
}