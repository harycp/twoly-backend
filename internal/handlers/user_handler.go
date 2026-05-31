package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harycp/twoly-backend/internal/services"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService}
}

func (h *UserHandler) UpdatePresence(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	err := h.userService.UpdateLastSeen(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update presence"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Presence updated successfully",
	})
}