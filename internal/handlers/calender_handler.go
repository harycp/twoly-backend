package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harycp/twoly-backend/internal/dto"
	"github.com/harycp/twoly-backend/internal/services"
)

type CalendarHandler struct {
	calendarService services.CalendarService
}

func NewCalendarHandler(calendarService services.CalendarService) *CalendarHandler {
	return &CalendarHandler{calendarService}
}

func (h *CalendarHandler) CreateCustomEvent(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req dto.CreateCustomEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "error": err.Error()})
		return
	}

	res, err := h.calendarService.CreateCustomEvent(userID.(string), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Custom event created successfully",
		"data":    res,
	})
}

func (h *CalendarHandler) GetEvents(c *gin.Context) {
	userID, _ := c.Get("userID")

	start := c.Query("start")
	end := c.Query("end")

	if start == "" || end == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Both 'start' and 'end' query parameters are required (YYYY-MM-DD)"})
		return
	}

	res, err := h.calendarService.GetCalendarEvents(userID.(string), start, end)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Calendar events fetched successfully",
		"data":    res,
	})
}

func (h *CalendarHandler) UpdateCustomEvent(c *gin.Context) {
	userID, _ := c.Get("userID")
	eventID := c.Param("id")

	var req dto.UpdateCustomEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "error": err.Error()})
		return
	}

	res, err := h.calendarService.UpdateCustomEvent(userID.(string), eventID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Custom event updated successfully",
		"data":    res,
	})
}

func (h *CalendarHandler) DeleteCustomEvent(c *gin.Context) {
	userID, _ := c.Get("userID")
	eventID := c.Param("id")

	err := h.calendarService.DeleteCustomEvent(userID.(string), eventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Custom event deleted successfully",
	})
}