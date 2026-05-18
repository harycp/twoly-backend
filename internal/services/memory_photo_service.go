package services

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/harycp/twoly-backend/internal/dto"
	"github.com/harycp/twoly-backend/internal/models"
	"github.com/harycp/twoly-backend/internal/repositories"
)

type MemoryPhotoService interface {
	UploadPhotos(userID string, memoryID string, fileHeaders []*multipart.FileHeader, captions []string) ([]dto.MemoryPhotoResponse, error)
	GetPhotosByMemoryID(userID string, memoryID string) ([]dto.MemoryPhotoResponse, error)
	GetGalleryPhotos(userID string) ([]dto.GalleryPhotoResponse, error)
	DeletePhoto(userID string, memoryID string, photoID string) error
}

type memoryPhotoService struct {
	photoRepo     repositories.MemoryPhotoRepository
	memoryRepo    repositories.MemoryRepository
	coupleRepo    repositories.CoupleRepository
	cloudinarySvc CloudinaryService
}

func NewMemoryPhotoService(
	photoRepo repositories.MemoryPhotoRepository,
	memoryRepo repositories.MemoryRepository,
	coupleRepo repositories.CoupleRepository,
	cloudinarySvc CloudinaryService,
) MemoryPhotoService {
	return &memoryPhotoService{photoRepo, memoryRepo, coupleRepo, cloudinarySvc}
}

// verifyMemoryAccess ensures the user belongs to the couple that owns the memory
func (s *memoryPhotoService) verifyMemoryAccess(userID string, memoryID string) (*models.Memory, error) {
	couple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil || couple == nil || couple.Status != "active" {
		return nil, errors.New("you do not have an active couple")
	}

	memory, err := s.memoryRepo.FindByID(memoryID, couple.ID.String())
	if err != nil || memory == nil {
		return nil, errors.New("memory not found or you do not have permission")
	}

	return memory, nil
}

func (s *memoryPhotoService) UploadPhotos(userID string, memoryID string, fileHeaders []*multipart.FileHeader, captions []string) ([]dto.MemoryPhotoResponse, error) {
	memory, err := s.verifyMemoryAccess(userID, memoryID)
	if err != nil {
		return nil, err
	}

	userUUID, _ := uuid.Parse(userID)
	var responses []dto.MemoryPhotoResponse

	for i, fileHeader := range fileHeaders {
		// Validasi tipe file
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
			continue // Skip file yang tidak valid
		}

		file, err := fileHeader.Open()
		if err != nil {
			continue
		}

		// Upload ke Cloudinary
		uniqueFilename := uuid.New().String()
		secureURL, publicID, err := s.cloudinarySvc.UploadImage(file, uniqueFilename)
		file.Close() // Pastikan file ditutup setelah diupload

		if err != nil {
			continue // Skip jika gagal upload ke cloud
		}

		// Atur caption (jika caption yang dikirim lebih sedikit dari jumlah foto, gunakan caption pertama/kosong)
		caption := ""
		if len(captions) > i {
			caption = captions[i]
		} else if len(captions) > 0 {
			caption = captions[0]
		}

		// Simpan ke Database
		photo := models.MemoryPhoto{
			MemoryID:           memory.ID,
			UploadedBy:         userUUID,
			PhotoURL:           secureURL,
			CloudinaryPublicID: publicID,
			Caption:            caption,
		}

		if err := s.photoRepo.CreatePhoto(&photo); err == nil {
			responses = append(responses, s.mapToResponse(photo))
		}
	}

	if len(responses) == 0 {
		return nil, errors.New("failed to upload any valid photos")
	}

	return responses, nil
}

func (s *memoryPhotoService) GetPhotosByMemoryID(userID string, memoryID string) ([]dto.MemoryPhotoResponse, error) {
	_, err := s.verifyMemoryAccess(userID, memoryID)
	if err != nil {
		return nil, err
	}

	photos, err := s.photoRepo.FindByMemoryID(memoryID)
	if err != nil {
		return nil, err
	}

	var response []dto.MemoryPhotoResponse
	for _, p := range photos {
		response = append(response, s.mapToResponse(p))
	}

	if response == nil {
		response = []dto.MemoryPhotoResponse{}
	}

	return response, nil
}

func (s *memoryPhotoService) GetGalleryPhotos(userID string) ([]dto.GalleryPhotoResponse, error) {
	couple, err := s.coupleRepo.FindByUserID(userID)
	if err != nil || couple == nil || couple.Status != "active" {
		return nil, errors.New("you do not have an active couple")
	}

	photos, err := s.photoRepo.FindAllByCoupleID(couple.ID.String())
	if err != nil {
		return nil, err
	}

	response := make([]dto.GalleryPhotoResponse, 0, len(photos))
	for _, p := range photos {
		response = append(response, s.mapToGalleryResponse(p))
	}

	return response, nil
}

func (s *memoryPhotoService) DeletePhoto(userID string, memoryID string, photoID string) error {
	_, err := s.verifyMemoryAccess(userID, memoryID)
	if err != nil {
		return err
	}

	photo, err := s.photoRepo.FindByID(photoID)
	if err != nil || photo.MemoryID.String() != memoryID {
		return errors.New("photo not found in this memory")
	}

	// Delete from Cloudinary
	if photo.CloudinaryPublicID != "" {
		_ = s.cloudinarySvc.DeleteImage(photo.CloudinaryPublicID)
	}

	// Delete from Database
	return s.photoRepo.DeletePhoto(photo)
}

func (s *memoryPhotoService) mapToResponse(p models.MemoryPhoto) dto.MemoryPhotoResponse {
	return dto.MemoryPhotoResponse{
		ID:                 p.ID.String(),
		MemoryID:           p.MemoryID.String(),
		UploadedBy:         p.UploadedBy.String(),
		PhotoURL:           p.PhotoURL,
		CloudinaryPublicID: p.CloudinaryPublicID,
		Caption:            p.Caption,
		CreatedAt:          p.CreatedAt,
	}
}

func (s *memoryPhotoService) mapToGalleryResponse(p models.MemoryPhoto) dto.GalleryPhotoResponse {
	return dto.GalleryPhotoResponse{
		ID:                 p.ID.String(),
		MemoryID:           p.MemoryID.String(),
		UploadedBy:         p.UploadedBy.String(),
		PhotoURL:           p.PhotoURL,
		CloudinaryPublicID: p.CloudinaryPublicID,
		Caption:            p.Caption,
		CreatedAt:          p.CreatedAt,
		Memory: dto.GalleryMemoryDetail{
			ID:           p.Memory.ID.String(),
			Title:        p.Memory.Title,
			Description:  p.Memory.Description,
			MemoryDate:   p.Memory.MemoryDate,
			LocationName: p.Memory.LocationName,
			Latitude:     p.Memory.Latitude,
			Longitude:    p.Memory.Longitude,
			Mood:         p.Memory.Mood,
			Tags:         []string(p.Memory.Tags),
			CreatedAt:    p.Memory.CreatedAt,
		},
	}
}