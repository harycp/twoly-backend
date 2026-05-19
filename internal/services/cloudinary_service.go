package services

import (
	"context"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/harycp/twoly-backend/internal/config"
)

type CloudinaryService interface {
	UploadImage(file multipart.File, filename string, folder string) (string, string, error)
	DeleteImage(publicID string) error
}

type cloudinaryService struct{}

func NewCloudinaryService() CloudinaryService {
	return &cloudinaryService{}
}

func (s *cloudinaryService) UploadImage(file multipart.File, filename string, folder string) (string, string, error) {
	cld := config.GetCloudinary()
	ctx := context.Background()
	if folder == "" {
		folder = os.Getenv("CLOUDINARY_FOLDER")
	}

	resp, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:   folder,
		PublicID: filename,
	})

	if err != nil {
		return "", "", err
	}

	return resp.SecureURL, resp.PublicID, nil
}

func (s *cloudinaryService) DeleteImage(publicID string) error {
	cld := config.GetCloudinary()
	ctx := context.Background()

	_, err := cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})

	return err
}