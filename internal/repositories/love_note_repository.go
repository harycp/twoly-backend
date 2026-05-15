package repositories

import (
	"github.com/harycp/twoly-backend/internal/models"
	"gorm.io/gorm"
)

type LoveNoteRepository interface {
	CreateLoveNote(note *models.LoveNote) error
	FindAllByCoupleID(coupleID string) ([]models.LoveNote, error)
	FindByID(id string, coupleID string) (*models.LoveNote, error)
	UpdateLoveNote(note *models.LoveNote) error
	DeleteLoveNote(note *models.LoveNote) error
}

type loveNoteRepository struct {
	db *gorm.DB
}

func NewLoveNoteRepository(db *gorm.DB) LoveNoteRepository {
	return &loveNoteRepository{db}
}

func (r *loveNoteRepository) CreateLoveNote(note *models.LoveNote) error {
	return r.db.Create(note).Error
}

func (r *loveNoteRepository) FindAllByCoupleID(coupleID string) ([]models.LoveNote, error) {
	var notes []models.LoveNote
	err := r.db.Where("couple_id = ?", coupleID).Order("created_at DESC").Find(&notes).Error
	return notes, err
}

func (r *loveNoteRepository) FindByID(id string, coupleID string) (*models.LoveNote, error) {
	var note models.LoveNote
	err := r.db.Where("id = ? AND couple_id = ?", id, coupleID).First(&note).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *loveNoteRepository) UpdateLoveNote(note *models.LoveNote) error {
	return r.db.Save(note).Error
}

func (r *loveNoteRepository) DeleteLoveNote(note *models.LoveNote) error {
	return r.db.Delete(note).Error
}