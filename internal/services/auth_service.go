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
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo}
}

func (s *authService) Register(req dto.RegisterRequest) (dto.AuthResponse, error) {
	// 1. Cek apakah email sudah terdaftar
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil && existingUser.Email != "" {
		return dto.AuthResponse{}, errors.New("email sudah terdaftar")
	}

	// 2. Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	// 3. Siapkan model User
	user := models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	// 4. Simpan ke database
	err = s.userRepo.CreateUser(&user)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	// 5. Generate JWT Token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{
		AccessToken: token,
		User: dto.UserResponse{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

func (s *authService) Login(req dto.LoginRequest) (dto.AuthResponse, error) {
	// 1. Cari user berdasarkan email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return dto.AuthResponse{}, errors.New("email atau password salah")
	}

	// 2. Bandingkan password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return dto.AuthResponse{}, errors.New("email atau password salah")
	}

	// 3. Generate JWT Token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{
		AccessToken: token,
		User: dto.UserResponse{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}