package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harycp/twoly-backend/internal/dto"
	"github.com/harycp/twoly-backend/internal/services"
)

type DatePlanHandler struct {
	datePlanService services.DatePlanService
}

func NewDatePlanHandler(datePlanService services.DatePlanService) *DatePlanHandler {
	return &DatePlanHandler{datePlanService}
}

func (h *DatePlanHandler) CreateDatePlan(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req dto.CreateDatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "error": err.Error()})
		return
	}

	res, err := h.datePlanService.CreateDatePlan(userID.(string), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Date plan created successfully",
		"data":    res,
	})
}

func (h *DatePlanHandler) GetAllDatePlans(c *gin.Context) {
	userID, _ := c.Get("userID")
	status := c.Query("status") // Allow filtering by ?status=planned

	res, err := h.datePlanService.GetAllDatePlans(userID.(string), status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Date plans fetched successfully",
		"data":    res,
	})
}

func (h *DatePlanHandler) GetDatePlanDetail(c *gin.Context) {
	userID, _ := c.Get("userID")
	planID := c.Param("id")

	res, err := h.datePlanService.GetDatePlanDetail(userID.(string), planID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Date plan fetched successfully",
		"data":    res,
	})
}

func (h *DatePlanHandler) UpdateDatePlan(c *gin.Context) {
	userID, _ := c.Get("userID")
	planID := c.Param("id")

	var req dto.UpdateDatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "error": err.Error()})
		return
	}

	res, err := h.datePlanService.UpdateDatePlan(userID.(string), planID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Date plan updated successfully",
		"data":    res,
	})
}

func (h *DatePlanHandler) DeleteDatePlan(c *gin.Context) {
	userID, _ := c.Get("userID")
	planID := c.Param("id")

	err := h.datePlanService.DeleteDatePlan(userID.(string), planID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Date plan deleted successfully",
	})
}

func (h *DatePlanHandler) UpdateStatus(c *gin.Context) {
	userID, _ := c.Get("userID")
	planID := c.Param("id")

	var req dto.UpdateDatePlanStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "error": err.Error()})
		return
	}

	res, err := h.datePlanService.UpdateDatePlanStatus(userID.(string), planID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Date plan status updated successfully",
		"data":    res,
	})
}

func (h *DatePlanHandler) ConvertToMemory(c *gin.Context) {
	userID, _ := c.Get("userID")
	planID := c.Param("id")

	res, err := h.datePlanService.ConvertToMemory(userID.(string), planID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Date plan successfully converted to memory",
		"data":    res, // Mengembalikan data memory yang baru saja dibuat
	})
}

func (h *DatePlanHandler) UpdateChecklistItem(c *gin.Context) {
	userID, _ := c.Get("userID")
	planID := c.Param("id")
	checklistID := c.Param("checklistId")

	var req dto.UpdateChecklistItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "error": err.Error()})
		return
	}

	res, err := h.datePlanService.UpdateChecklistItem(userID.(string), planID, checklistID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Checklist item updated successfully",
		"data":    res,
	})
}

func (h *DatePlanHandler) AddChecklistItem(c *gin.Context) {
	userID, _ := c.Get("userID")
	planID := c.Param("id")

	var req dto.AddChecklistItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "error": err.Error()})
		return
	}

	res, err := h.datePlanService.AddChecklistItem(userID.(string), planID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Checklist item added successfully",
		"data":    res,
	})
}

func (h *DatePlanHandler) DeleteChecklistItem(c *gin.Context) {
	userID, _ := c.Get("userID")
	planID := c.Param("id")
	checklistID := c.Param("checklistId")

	err := h.datePlanService.DeleteChecklistItem(userID.(string), planID, checklistID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Checklist item deleted successfully",
	})
}