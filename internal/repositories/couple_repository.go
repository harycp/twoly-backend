package repositories

import (
	"github.com/harycp/twoly-backend/internal/models"
	"gorm.io/gorm"
)

type CoupleRepository interface {
	CreateCouple(couple *models.Couple) error
	FindByInviteCode(code string) (*models.Couple, error)
	FindByUserID(userID string) (*models.Couple, error)
	UpdateCouple(couple *models.Couple) error
}

type coupleRepository struct {
	db *gorm.DB
}

func NewCoupleRepository(db *gorm.DB) CoupleRepository {
	return &coupleRepository{db}
}

func (r *coupleRepository) CreateCouple(couple *models.Couple) error {
	return r.db.Create(couple).Error
}

func (r *coupleRepository) FindByInviteCode(code string) (*models.Couple, error) {
	var couple models.Couple
	err := r.db.Preload("UserOne").Preload("UserTwo").Where("invite_code = ?", code).First(&couple).Error
	if err != nil {
		return nil, err
	}
	return &couple, nil
}

func (r *coupleRepository) FindByUserID(userID string) (*models.Couple, error) {
	var couple models.Couple
	err := r.db.Preload("UserOne").Preload("UserTwo").
		Where("user_one_id = ? OR user_two_id = ?", userID, userID).
		First(&couple).Error
	if err != nil {
		return nil, err
	}
	return &couple, nil
}

func (r *coupleRepository) UpdateCouple(couple *models.Couple) error {
	return r.db.Save(couple).Error
}