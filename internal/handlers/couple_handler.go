package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harycp/twoly-backend/internal/dto"
	"github.com/harycp/twoly-backend/internal/services"
)

type CoupleHandler struct {
	coupleService services.CoupleService
}

func NewCoupleHandler(coupleService services.CoupleService) *CoupleHandler {
	return &CoupleHandler{coupleService}
}

func (h *CoupleHandler) CreateInvite(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req dto.CreateInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil && err.Error() != "EOF" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "error": err.Error()})
		return
	}

	res, err := h.coupleService.CreateInvite(userID.(string), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Invite created",
		"data":    res,
	})
}

func (h *CoupleHandler) JoinCouple(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req dto.JoinCoupleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "error": err.Error()})
		return
	}

	err := h.coupleService.JoinCouple(userID.(string), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Couple joined successfully",
	})
}

func (h *CoupleHandler) GetMyCouple(c *gin.Context) {
	userID, _ := c.Get("userID")

	res, err := h.coupleService.GetMyCouple(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Couple profile fetched successfully",
		"data":    res,
	})
}