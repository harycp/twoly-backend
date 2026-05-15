package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harycp/twoly-backend/internal/services"
)

type MemoryPhotoHandler struct {
	photoService services.MemoryPhotoService
}

func NewMemoryPhotoHandler(photoService services.MemoryPhotoService) *MemoryPhotoHandler {
	return &MemoryPhotoHandler{photoService}
}

func (h *MemoryPhotoHandler) UploadPhotos(c *gin.Context) {
	userID, _ := c.Get("userID")
	memoryID := c.Param("id")

	// Limit upload size (e.g., 20MB total untuk banyak file)
	err := c.Request.ParseMultipartForm(20 << 20)
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
		c.JSON(http.StatusBadRequest, gin.H{"message": "At least one image file is required"})
		return
	}

	captions := form.Value["captions"]

	res, err := h.photoService.UploadPhotos(userID.(string), memoryID, files, captions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Photos uploaded successfully",
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
		"message": "Photos fetched successfully",
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