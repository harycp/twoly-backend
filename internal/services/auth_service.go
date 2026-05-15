package services

import (
	"errors"

	"github.com/harycp/twoly-backend/internal/dto"
	"github.com/harycp/twoly-backend/internal/models"
	"github.com/harycp/twoly-backend/internal/repositories"
	"github.com/harycp/twoly-backend/internal/utils"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (dto.AuthResponse, error)
	Login(req dto.LoginRequest) (dto.AuthResponse, error)
	GetMe(userID string) (dto.UserResponse, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo}
}

func (s *authService) Register(req dto.RegisterRequest) (dto.AuthResponse, error) {
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil && existingUser.Email != "" {
		return dto.AuthResponse{}, errors.New("email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	user := models.User{
		Name:         req.Name,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	err = s.userRepo.CreateUser(&user)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{
		AccessToken: token,
		User: dto.UserResponse{
			ID:       user.ID.String(),
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
		},
	}, nil
}

func (s *authService) Login(req dto.LoginRequest) (dto.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return dto.AuthResponse{}, errors.New("email atau password salah")
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return dto.AuthResponse{}, errors.New("email atau password salah")
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{
		AccessToken: token,
		User: dto.UserResponse{
			ID:       user.ID.String(),
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
		},
	}, nil
}

func (s *authService) GetMe(userID string) (dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return dto.UserResponse{}, errors.New("user tidak ditemukan")
	}

	return dto.UserResponse{
		ID:       user.ID.String(),
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}