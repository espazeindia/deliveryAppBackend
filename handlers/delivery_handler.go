package handlers

import (
	"deliveryAppBackend/domain/entities"
	"deliveryAppBackend/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeliveryHandler struct {
	deliveryUseCase *usecase.DeliveryUseCase
}

func NewDeliveryHandler(deliveryUseCase *usecase.DeliveryUseCase) *DeliveryHandler {
	return &DeliveryHandler{
		deliveryUseCase: deliveryUseCase,
	}
}

func (h *DeliveryHandler) GetActiveOrders(c *gin.Context) {
	partnerID := c.GetString("partnerId")
	
	response, err := h.deliveryUseCase.GetActiveOrders(partnerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *DeliveryHandler) GetOrderHistory(c *gin.Context) {
	partnerID := c.GetString("partnerId")
	
	var req entities.GetOrderHistoryRequest
	req.Limit = 20 // default
	req.Offset = 0  // default

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response, err := h.deliveryUseCase.GetOrderHistory(partnerID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *DeliveryHandler) GetOrderDetails(c *gin.Context) {
	deliveryID := c.Param("id")
	
	response, err := h.deliveryUseCase.GetOrderDetails(deliveryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *DeliveryHandler) AcceptOrder(c *gin.Context) {
	deliveryID := c.Param("id")
	partnerID := c.GetString("partnerId")
	
	response, err := h.deliveryUseCase.AcceptOrder(deliveryID, partnerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *DeliveryHandler) UpdateOrderStatus(c *gin.Context) {
	deliveryID := c.Param("id")
	partnerID := c.GetString("partnerId")
	
	var req entities.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response, err := h.deliveryUseCase.UpdateOrderStatus(deliveryID, partnerID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *DeliveryHandler) CompleteDelivery(c *gin.Context) {
	deliveryID := c.Param("id")
	partnerID := c.GetString("partnerId")
	
	var req entities.CompleteDeliveryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response, err := h.deliveryUseCase.CompleteDelivery(deliveryID, partnerID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

