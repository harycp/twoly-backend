package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harycp/twoly-backend/internal/dto"
	"github.com/harycp/twoly-backend/internal/services"
)

type MemoryHandler struct {
	memoryService services.MemoryService
}

func NewMemoryHandler(memoryService services.MemoryService) *MemoryHandler {
	return &MemoryHandler{memoryService}
}

func (h *MemoryHandler) CreateMemory(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req dto.CreateMemoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "error": err.Error()})
		return
	}

	res, err := h.memoryService.CreateMemory(userID.(string), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Memory created successfully",
		"data":    res,
	})
}

func (h *MemoryHandler) GetAllMemories(c *gin.Context) {
	userID, _ := c.Get("userID")
	month := c.Query("month") // Menangkap query parameter ?month=YYYY-MM

	res, err := h.memoryService.GetAllMemories(userID.(string), month)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Memories fetched successfully",
		"data":    res,
	})
}

func (h *MemoryHandler) GetMemoryDetail(c *gin.Context) {
	userID, _ := c.Get("userID")
	memoryID := c.Param("id")

	res, err := h.memoryService.GetMemoryDetail(userID.(string), memoryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Memory fetched successfully",
		"data":    res,
	})
}

func (h *MemoryHandler) UpdateMemory(c *gin.Context) {
	userID, _ := c.Get("userID")
	memoryID := c.Param("id")

	var req dto.UpdateMemoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "error": err.Error()})
		return
	}

	res, err := h.memoryService.UpdateMemory(userID.(string), memoryID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Memory updated successfully",
		"data":    res,
	})
}

func (h *MemoryHandler) DeleteMemory(c *gin.Context) {
	userID, _ := c.Get("userID")
	memoryID := c.Param("id")

	err := h.memoryService.DeleteMemory(userID.(string), memoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Memory deleted successfully",
	})
}