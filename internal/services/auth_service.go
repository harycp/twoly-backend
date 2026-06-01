package services

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/harycp/twoly-backend/internal/dto"
	"github.com/harycp/twoly-backend/internal/models"
	"github.com/harycp/twoly-backend/internal/repositories"
	"github.com/harycp/twoly-backend/internal/utils"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (dto.AuthResponse, error)
	Login(req dto.LoginRequest) (dto.AuthResponse, error)
	GoogleLogin(req dto.GoogleLoginRequest) (dto.AuthResponse, error)
	GetMe(userID string) (dto.UserResponse, error)
	UpdateProfile(userID string, req dto.UpdateProfileRequest) (dto.UserResponse, error)
	ForgotPassword(req dto.ForgotPasswordRequest) error
	VerifyOTP(req dto.VerifyOTPRequest) error
	ResetPassword(req dto.ResetPasswordRequest) error
}

type authService struct {
	userRepo      repositories.UserRepository
	cloudinarySvc CloudinaryService
	emailSvc      EmailService
}

func NewAuthService(userRepo repositories.UserRepository, cloudinarySvc CloudinaryService, emailSvc EmailService) AuthService {
	return &authService{userRepo, cloudinarySvc, emailSvc}
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
	user, err := s.userRepo.FindByEmailOrUsername(req.EmailOrUsername)
	if err != nil {
		return dto.AuthResponse{}, errors.New("invalid email, username, or password")
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return dto.AuthResponse{}, errors.New("invalid email, username, or password")
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

func (s *authService) GoogleLogin(req dto.GoogleLoginRequest) (dto.AuthResponse, error) {
	resp, err := http.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + req.IDToken)
	if err != nil || resp.StatusCode != http.StatusOK {
		return dto.AuthResponse{}, errors.New("invalid google token")
	}
	defer resp.Body.Close()

	var googleClaims struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&googleClaims); err != nil {
		return dto.AuthResponse{}, errors.New("failed to parse google token")
	}

	user, err := s.userRepo.FindByEmail(googleClaims.Email)
	if err != nil || user.ID == uuid.Nil {
		baseUsername := strings.ReplaceAll(strings.ToLower(googleClaims.Name), " ", "")
		randomString := utils.GenerateInviteCode(4)
		newUsername := fmt.Sprintf("%s_%s", baseUsername, randomString)
		randomPassword, _ := utils.HashPassword(uuid.NewString())

		newUser := models.User{
			Name:         googleClaims.Name,
			Username:     newUsername,
			Email:        googleClaims.Email,
			PasswordHash: randomPassword,
			AvatarURL:    &googleClaims.Picture,
		}

		if err := s.userRepo.CreateUser(&newUser); err != nil {
			return dto.AuthResponse{}, errors.New("failed to create user via google")
		}
		user = &newUser
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{
		AccessToken: token,
		User: dto.UserResponse{
			ID:        user.ID.String(),
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			AvatarURL: user.AvatarURL,
		},
	}, nil
}

func (s *authService) GetMe(userID string) (dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return dto.UserResponse{}, errors.New("user not found")
	}

	return dto.UserResponse{
		ID:        user.ID.String(),
		Username:  user.Username,
		Name:      user.Name,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
	}, nil
}

func (s *authService) UpdateProfile(userID string, req dto.UpdateProfileRequest) (dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return dto.UserResponse{}, errors.New("user not found")
	}

	if req.Username != "" && req.Username != user.Username {
		existingUsername, _ := s.userRepo.FindByUsername(req.Username)
		if existingUsername != nil && existingUsername.Username != "" {
			return dto.UserResponse{}, errors.New("username is already taken by another user")
		}
		user.Username = req.Username
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Avatar != nil {
		file, err := req.Avatar.Open()
		if err != nil {
			return dto.UserResponse{}, errors.New("failed to read avatar file")
		}

		oldAvatarPublicID := user.AvatarPublicID
		avatarPublicID := user.ID.String() + "_avatar_" + uuid.NewString()
		secureURL, newPublicID, err := s.cloudinarySvc.UploadImage(file, avatarPublicID, "twoly/avatar")
		file.Close()
		if err != nil {
			return dto.UserResponse{}, errors.New("failed to upload avatar")
		}

		if oldAvatarPublicID != nil && *oldAvatarPublicID != "" {
			_ = s.cloudinarySvc.DeleteImage(*oldAvatarPublicID)
		}

		user.AvatarURL = &secureURL
		user.AvatarPublicID = &newPublicID
	}

	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
		user.AvatarPublicID = nil
	}

	if err := s.userRepo.UpdateUser(user); err != nil {
		return dto.UserResponse{}, errors.New("failed to update profile")
	}

	return dto.UserResponse{
		ID:        user.ID.String(),
		Username:  user.Username,
		Name:      user.Name,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
	}, nil
}

// ------------------------------------------------------------------
// FORGOT PASSWORD LOGIC
// ------------------------------------------------------------------

func generateNumericOTP(length int) string {
	const charset = "0123456789"
	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

func (s *authService) ForgotPassword(req dto.ForgotPasswordRequest) error {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil || user.ID == uuid.Nil {
		return errors.New("no account found with this email")
	}

	if user.ResetOTPExpiresAt != nil && time.Now().Before(*user.ResetOTPExpiresAt) {
		remaining := time.Until(*user.ResetOTPExpiresAt)
		minutes := int(remaining.Minutes())
		
		if minutes > 0 {
			return fmt.Errorf("an active OTP was already sent. please wait %d minutes", minutes)
		}
		return errors.New("an active OTP was already sent. please wait a few seconds")
	}

	otp := generateNumericOTP(6)
	expiry := time.Now().Add(15 * time.Minute)

	user.ResetOTP = &otp
	user.ResetOTPExpiresAt = &expiry

	if err := s.userRepo.UpdateUser(user); err != nil {
		return errors.New("failed to generate OTP")
	}

	err = s.emailSvc.SendPasswordResetOTP(user.Email, user.Name, otp)
	if err != nil {
		return errors.New("failed to send email. please try again later")
	}

	return nil
}

func (s *authService) VerifyOTP(req dto.VerifyOTPRequest) error {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil || user.ID == uuid.Nil {
		return errors.New("invalid email or OTP")
	}

	if user.ResetOTP == nil || *user.ResetOTP != req.OTP {
		return errors.New("invalid OTP code")
	}

	if user.ResetOTPExpiresAt == nil || time.Now().After(*user.ResetOTPExpiresAt) {
		return errors.New("OTP code has expired")
	}

	return nil
}

func (s *authService) ResetPassword(req dto.ResetPasswordRequest) error {
	err := s.VerifyOTP(dto.VerifyOTPRequest{Email: req.Email, OTP: req.OTP})
	if err != nil {
		return err
	}

	user, _ := s.userRepo.FindByEmail(req.Email)

	hashed, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.PasswordHash = hashed
	user.ResetOTP = nil
	user.ResetOTPExpiresAt = nil

	if err := s.userRepo.UpdateUser(user); err != nil {
		return errors.New("failed to update password")
	}

	return nil
}