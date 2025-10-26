package handlers

import (
	"deliveryAppBackend/domain/entities"
	"deliveryAppBackend/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	profileUseCase *usecase.ProfileUseCase
}

func NewProfileHandler(profileUseCase *usecase.ProfileUseCase) *ProfileHandler {
	return &ProfileHandler{
		profileUseCase: profileUseCase,
	}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	partnerID := c.GetString("partnerId")
	
	profile, err := h.profileUseCase.GetProfile(partnerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"profile": profile,
	})
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	partnerID := c.GetString("partnerId")
	
	var req entities.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response, err := h.profileUseCase.UpdateProfile(partnerID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProfileHandler) UpdateLocation(c *gin.Context) {
	partnerID := c.GetString("partnerId")
	
	var req entities.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response, err := h.profileUseCase.UpdateLocation(partnerID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProfileHandler) ToggleAvailability(c *gin.Context) {
	partnerID := c.GetString("partnerId")
	
	var req entities.ToggleAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response, err := h.profileUseCase.ToggleAvailability(partnerID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

