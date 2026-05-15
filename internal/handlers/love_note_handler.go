package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harycp/twoly-backend/internal/dto"
	"github.com/harycp/twoly-backend/internal/services"
)

type LoveNoteHandler struct {
	loveNoteService services.LoveNoteService
}

func NewLoveNoteHandler(loveNoteService services.LoveNoteService) *LoveNoteHandler {
	return &LoveNoteHandler{loveNoteService}
}

func (h *LoveNoteHandler) CreateLoveNote(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req dto.CreateLoveNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "error": err.Error()})
		return
	}

	res, err := h.loveNoteService.CreateLoveNote(userID.(string), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Love note sent successfully",
		"data":    res,
	})
}

func (h *LoveNoteHandler) GetLoveNotes(c *gin.Context) {
	userID, _ := c.Get("userID")

	res, err := h.loveNoteService.GetLoveNotes(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Love notes fetched successfully",
		"data":    res,
	})
}

func (h *LoveNoteHandler) OpenLoveNote(c *gin.Context) {
	userID, _ := c.Get("userID")
	noteID := c.Param("id")

	res, err := h.loveNoteService.OpenLoveNote(userID.(string), noteID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Love note opened successfully",
		"data":    res,
	})
}

func (h *LoveNoteHandler) DeleteLoveNote(c *gin.Context) {
	userID, _ := c.Get("userID")
	noteID := c.Param("id")

	err := h.loveNoteService.DeleteLoveNote(userID.(string), noteID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Love note deleted successfully",
	})
}