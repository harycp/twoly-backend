package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/harycp/twoly-backend/internal/dto"
	"github.com/harycp/twoly-backend/internal/models"
	"github.com/harycp/twoly-backend/internal/repositories"
)

type LoveNoteService interface {
	CreateLoveNote(userID string, req dto.CreateLoveNoteRequest) (dto.LoveNoteResponse, error)
	GetLoveNotes(userID string) ([]dto.LoveNoteResponse, error)
	OpenLoveNote(userID string, noteID string) (dto.LoveNoteResponse, error)
	DeleteLoveNote(userID string, noteID string) error
}

type loveNoteService struct {
	loveNoteRepo repositories.LoveNoteRepository
	coupleRepo   repositories.CoupleRepository
}

func NewLoveNoteService(loveNoteRepo repositories.LoveNoteRepository, coupleRepo repositories.CoupleRepository) LoveNoteService {
	return &loveNoteService{loveNoteRepo, coupleRepo}
}

func (s *loveNoteService) getActiveCoupleAndPartner(userID string) (*models.Couple, uuid.UUID, error) {
	couple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil || couple == nil {
		return nil, uuid.Nil, errors.New("you do not have an active couple")
	}

	if couple.Status != "active" || couple.UserTwoID == nil {
		return nil, uuid.Nil, errors.New("your partner hasn't joined yet")
	}

	var partnerID uuid.UUID
	if couple.UserOneID.String() == userID {
		partnerID = *couple.UserTwoID
	} else {
		partnerID = couple.UserOneID
	}

	return couple, partnerID, nil
}

func (s *loveNoteService) CreateLoveNote(userID string, req dto.CreateLoveNoteRequest) (dto.LoveNoteResponse, error) {
	couple, partnerID, err := s.getActiveCoupleAndPartner(userID)
	if err != nil {
		return dto.LoveNoteResponse{}, err
	}

	userUUID, _ := uuid.Parse(userID)
	
	var unlockAt *time.Time
	if req.UnlockAt != "" {
		t, err := time.Parse(time.RFC3339, req.UnlockAt)
		if err == nil {
			unlockAt = &t
		} else {
			return dto.LoveNoteResponse{}, errors.New("invalid unlock_at date format, please use ISO 8601")
		}
	}

	note := models.LoveNote{
		CoupleID:   couple.ID,
		SenderID:   userUUID,
		ReceiverID: partnerID,
		Message:    req.Message,
		UnlockAt:   unlockAt,
		IsOpened:   false,
	}

	if err := s.loveNoteRepo.CreateLoveNote(&note); err != nil {
		return dto.LoveNoteResponse{}, errors.New("failed to create love note")
	}

	return s.mapToResponse(note, false), nil
}

func (s *loveNoteService) GetLoveNotes(userID string) ([]dto.LoveNoteResponse, error) {
	couple, _, err := s.getActiveCoupleAndPartner(userID)
	if err != nil {
		return nil, err
	}

	notes, err := s.loveNoteRepo.FindAllByCoupleID(couple.ID.String())
	if err != nil {
		return nil, err
	}

	var response []dto.LoveNoteResponse
	for _, n := range notes {
		isCensored := false
		if !n.IsOpened && n.UnlockAt != nil && n.UnlockAt.After(time.Now()) && n.ReceiverID.String() == userID {
			isCensored = true
		}
		response = append(response, s.mapToResponse(n, isCensored))
	}

	if response == nil {
		response = []dto.LoveNoteResponse{}
	}

	return response, nil
}

func (s *loveNoteService) OpenLoveNote(userID string, noteID string) (dto.LoveNoteResponse, error) {
	couple, _, err := s.getActiveCoupleAndPartner(userID)
	if err != nil {
		return dto.LoveNoteResponse{}, err
	}

	note, err := s.loveNoteRepo.FindByID(noteID, couple.ID.String())
	if err != nil || note == nil {
		return dto.LoveNoteResponse{}, errors.New("love note not found")
	}

	if note.ReceiverID.String() != userID {
		return dto.LoveNoteResponse{}, errors.New("you can only open love notes sent to you")
	}

	if note.IsOpened {
		return s.mapToResponse(*note, false), nil
	}

	if note.UnlockAt != nil && note.UnlockAt.After(time.Now()) {
		return dto.LoveNoteResponse{}, errors.New("this love note is still locked until " + note.UnlockAt.Format(time.RFC3339))
	}

	now := time.Now()
	note.IsOpened = true
	note.OpenedAt = &now

	if err := s.loveNoteRepo.UpdateLoveNote(note); err != nil {
		return dto.LoveNoteResponse{}, errors.New("failed to open love note")
	}

	return s.mapToResponse(*note, false), nil
}

func (s *loveNoteService) DeleteLoveNote(userID string, noteID string) error {
	couple, _, err := s.getActiveCoupleAndPartner(userID)
	if err != nil {
		return err
	}

	note, err := s.loveNoteRepo.FindByID(noteID, couple.ID.String())
	if err != nil || note == nil {
		return errors.New("love note not found")
	}

	if note.SenderID.String() != userID {
		return errors.New("you can only delete love notes that you sent")
	}

	return s.loveNoteRepo.DeleteLoveNote(note)
}

func (s *loveNoteService) mapToResponse(n models.LoveNote, isCensored bool) dto.LoveNoteResponse {
	msg := n.Message
	if isCensored {
		msg = "[This love note is locked. Wait until the time comes!]"
	}

	return dto.LoveNoteResponse{
		ID:         n.ID.String(),
		CoupleID:   n.CoupleID.String(),
		SenderID:   n.SenderID.String(),
		ReceiverID: n.ReceiverID.String(),
		Message:    msg,
		UnlockAt:   n.UnlockAt,
		IsOpened:   n.IsOpened,
		OpenedAt:   n.OpenedAt,
		CreatedAt:  n.CreatedAt,
	}
}