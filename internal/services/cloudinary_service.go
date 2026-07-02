package services

import (
	"context"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/harycp/twoly-backend/internal/config"
)

type CloudinaryService interface {
	UploadMedia(file multipart.File, filename string, folder string, resourceType string) (string, string, error)
	UploadImage(file multipart.File, filename string, folder string) (string, string, error)
	DeleteMedia(publicID string, resourceType string) error
	DeleteImage(publicID string) error
}

type cloudinaryService struct{}

func NewCloudinaryService() CloudinaryService {
	return &cloudinaryService{}
}

func uploadTransformation(resourceType string) string {
	switch resourceType {
	case "video":
		return "c_limit,w_1280,q_auto:eco,f_mp4"
	default:
		return "c_limit,w_1600,q_auto:eco,f_auto"
	}
}

func (s *cloudinaryService) UploadMedia(file multipart.File, filename string, folder string, resourceType string) (string, string, error) {
	cld := config.GetCloudinary()
	ctx := context.Background()
	if folder == "" {
		folder = os.Getenv("CLOUDINARY_FOLDER")
	}

	if resourceType == "" {
		resourceType = "auto"
	}

	resp, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:         folder,
		PublicID:       filename,
		ResourceType:   resourceType,
		Transformation: uploadTransformation(resourceType),
	})

	if err != nil {
		return "", "", err
	}

	return resp.SecureURL, resp.PublicID, nil
}

func (s *cloudinaryService) UploadImage(file multipart.File, filename string, folder string) (string, string, error) {
	return s.UploadMedia(file, filename, folder, "image")
}

func (s *cloudinaryService) DeleteMedia(publicID string, resourceType string) error {
	cld := config.GetCloudinary()
	ctx := context.Background()

	if resourceType == "" {
		resourceType = "image"
	}

	_, err := cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: resourceType,
	})

	return err
}

func (s *cloudinaryService) DeleteImage(publicID string) error {
	return s.DeleteMedia(publicID, "image")
}