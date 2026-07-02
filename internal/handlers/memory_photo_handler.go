package handlers

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/harycp/twoly-backend/internal/services"
)

type MemoryPhotoHandler struct {
	photoService services.MemoryPhotoService
}

func maxUploadSizeBytes() int64 {
	defaultSizeMB := int64(100)
	if value := os.Getenv("MAX_MEDIA_UPLOAD_SIZE_MB"); value != "" {
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil && parsed > 0 {
			defaultSizeMB = parsed
		}
	}

	return defaultSizeMB << 20
}

func NewMemoryPhotoHandler(photoService services.MemoryPhotoService) *MemoryPhotoHandler {
	return &MemoryPhotoHandler{photoService}
}

func (h *MemoryPhotoHandler) UploadPhotos(c *gin.Context) {
	userID, _ := c.Get("userID")
	memoryID := c.Param("id")

	err := c.Request.ParseMultipartForm(maxUploadSizeBytes())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Total file size exceeds the limit"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid form data"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "At least one image or video file is required"})
		return
	}

	captions := form.Value["captions"]

	res, err := h.photoService.UploadPhotos(userID.(string), memoryID, files, captions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Media uploaded successfully",
		"data":    res,
	})
}

func (h *MemoryPhotoHandler) GetPhotos(c *gin.Context) {
	userID, _ := c.Get("userID")
	memoryID := c.Param("id")

	res, err := h.photoService.GetPhotosByMemoryID(userID.(string), memoryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Media fetched successfully",
		"data":    res,
	})
}

func (h *MemoryPhotoHandler) GetGalleryPhotos(c *gin.Context) {
	userID, _ := c.Get("userID")

	res, err := h.photoService.GetGalleryPhotos(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Gallery photos fetched successfully",
		"data":    res,
	})
}

func (h *MemoryPhotoHandler) DeletePhoto(c *gin.Context) {
	userID, _ := c.Get("userID")
	memoryID := c.Param("id")
	photoID := c.Param("photoId")

	err := h.photoService.DeletePhoto(userID.(string), memoryID, photoID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Photo deleted successfully",
	})
}