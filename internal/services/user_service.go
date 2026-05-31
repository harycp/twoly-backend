package services

import (
	"time"

	"github.com/harycp/twoly-backend/internal/repositories"
)

type UserService interface {
	UpdateLastSeen(userID string) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo}
}

func (s *userService) UpdateLastSeen(userID string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	now := time.Now()
	user.LastSeen = &now
	
	return s.userRepo.UpdateUser(user)
}