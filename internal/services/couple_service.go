package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/harycp/twoly-backend/internal/dto"
	"github.com/harycp/twoly-backend/internal/models"
	"github.com/harycp/twoly-backend/internal/repositories"
	"github.com/harycp/twoly-backend/internal/utils"
)

type CoupleService interface {
	CreateInvite(userID string, req dto.CreateInviteRequest) (dto.CoupleResponse, error)
	JoinCouple(userID string, req dto.JoinCoupleRequest) error
	GetMyCouple(userID string) (dto.CoupleResponse, error)
	UpdateCoupleSettings(userID string, req dto.UpdateCoupleSettingsRequest) (dto.CoupleResponse, error)
}

type coupleService struct {
	coupleRepo repositories.CoupleRepository
}

func NewCoupleService(coupleRepo repositories.CoupleRepository) CoupleService {
	return &coupleService{coupleRepo}
}

func (s *coupleService) CreateInvite(userID string, req dto.CreateInviteRequest) (dto.CoupleResponse, error) {
	existingCouple, err := s.coupleRepo.FindByUserID(userID)
	if err == nil && existingCouple != nil {
		return dto.CoupleResponse{}, errors.New("user already in a couple or has a pending invite")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return dto.CoupleResponse{}, errors.New("invalid user ID")
	}

	var anniversary *time.Time
	if req.AnniversaryDate != nil {
		t, err := time.Parse("2006-01-02", *req.AnniversaryDate)
		if err == nil {
			anniversary = &t
		}
	}

	couple := models.Couple{
		UserOneID:       userUUID,
		InviteCode:      utils.GenerateInviteCode(6),
		AnniversaryDate: anniversary,
		Status:          "pending",
	}

	err = s.coupleRepo.CreateCouple(&couple)
	if err != nil {
		return dto.CoupleResponse{}, err
	}

	return dto.CoupleResponse{
		CoupleID:        couple.ID.String(),
		InviteCode:      couple.InviteCode,
		AnniversaryDate: couple.AnniversaryDate,
		Status:          couple.Status,
	}, nil
}

func (s *coupleService) JoinCouple(userID string, req dto.JoinCoupleRequest) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	existingCouple, err := s.coupleRepo.FindByUserID(userID)
	if err == nil && existingCouple != nil {
		return errors.New("you are already in a couple")
	}

	couple, err := s.coupleRepo.FindByInviteCode(req.InviteCode)
	if err != nil {
		return errors.New("invalid invite code")
	}

	if couple.Status != "pending" || couple.UserTwoID != nil {
		return errors.New("invite code is no longer valid")
	}

	if couple.UserOneID.String() == userID {
		return errors.New("you cannot join your own invite")
	}

	couple.UserTwoID = &userUUID
	couple.Status = "active"

	return s.coupleRepo.UpdateCouple(couple)
}

func (s *coupleService) GetMyCouple(userID string) (dto.CoupleResponse, error) {
	couple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil {
		return dto.CoupleResponse{}, errors.New("couple not found")
	}

	response := dto.CoupleResponse{
		CoupleID:        couple.ID.String(),
		InviteCode:      couple.InviteCode,
		AnniversaryDate: couple.AnniversaryDate,
		Status:          couple.Status,
	}

	var partner *models.User
	if couple.UserOneID.String() != userID && couple.UserOneID != uuid.Nil {
		partner = &couple.UserOne
	} else if couple.UserTwoID != nil && couple.UserTwoID.String() != userID {
		partner = couple.UserTwo
	}

	if partner != nil && partner.ID != uuid.Nil {
		response.Partner = &dto.UserResponse{
			ID:        partner.ID.String(),
			Name:      partner.Name,
			Username:  partner.Username,
			Email:     partner.Email,
			AvatarURL: partner.AvatarURL,
			LastSeen:  partner.LastSeen,
		}
	}

	return response, nil
}

func (s *coupleService) UpdateCoupleSettings(userID string, req dto.UpdateCoupleSettingsRequest) (dto.CoupleResponse, error) {
	couple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil || couple == nil {
		return dto.CoupleResponse{}, errors.New("couple not found")
	}

	if req.AnniversaryDate != nil {
		t, err := time.Parse("2006-01-02", *req.AnniversaryDate)
		if err == nil {
			couple.AnniversaryDate = &t
		} else {
			return dto.CoupleResponse{}, errors.New("invalid date format, use YYYY-MM-DD")
		}
	}

	if err := s.coupleRepo.UpdateCouple(couple); err != nil {
		return dto.CoupleResponse{}, errors.New("failed to update couple settings")
	}

	return s.GetMyCouple(userID)
}