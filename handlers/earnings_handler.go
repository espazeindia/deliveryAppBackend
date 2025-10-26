package handlers

import (
	"deliveryAppBackend/domain/entities"
	"deliveryAppBackend/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EarningsHandler struct {
	earningsUseCase *usecase.EarningsUseCase
}

func NewEarningsHandler(earningsUseCase *usecase.EarningsUseCase) *EarningsHandler {
	return &EarningsHandler{
		earningsUseCase: earningsUseCase,
	}
}

func (h *EarningsHandler) GetEarnings(c *gin.Context) {
	partnerID := c.GetString("partnerId")
	period := c.DefaultQuery("period", "week")
	
	response, err := h.earningsUseCase.GetEarnings(partnerID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *EarningsHandler) GetEarningsHistory(c *gin.Context) {
	partnerID := c.GetString("partnerId")
	
	var req entities.GetEarningsHistoryRequest
	req.Limit = 50 // default
	req.Offset = 0  // default

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response, err := h.earningsUseCase.GetEarningsHistory(partnerID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

