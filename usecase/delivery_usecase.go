package usecase

import (
	"deliveryAppBackend/domain/entities"
	"deliveryAppBackend/domain/repositories"
	"errors"
	"time"
)

type DeliveryUseCase struct {
	deliveryRepo repositories.DeliveryRepository
	partnerRepo  repositories.DeliveryPartnerRepository
	earningsRepo repositories.EarningsRepository
}

func NewDeliveryUseCase(
	deliveryRepo repositories.DeliveryRepository,
	partnerRepo repositories.DeliveryPartnerRepository,
	earningsRepo repositories.EarningsRepository,
) *DeliveryUseCase {
	return &DeliveryUseCase{
		deliveryRepo: deliveryRepo,
		partnerRepo:  partnerRepo,
		earningsRepo: earningsRepo,
	}
}

func (uc *DeliveryUseCase) GetActiveOrders(partnerID string) (*entities.GetActiveOrdersResponse, error) {
	deliveries, err := uc.deliveryRepo.GetActiveOrdersByPartner(partnerID)
	if err != nil {
		return &entities.GetActiveOrdersResponse{
			Success: false,
		}, err
	}

	orders := make([]entities.DeliveryListItem, 0, len(deliveries))
	for _, d := range deliveries {
		orders = append(orders, entities.DeliveryListItem{
			ID:          d.DeliveryID,
			OrderID:     d.OrderID,
			Status:      d.Status,
			Address:     d.DeliveryAddress,
			Amount:      d.OrderAmount,
			DeliveryFee: d.DeliveryFee,
			ItemsCount:  d.ItemsCount,
			Distance:    d.Distance,
			CreatedAt:   d.CreatedAt,
		})
	}

	return &entities.GetActiveOrdersResponse{
		Success: true,
		Orders:  orders,
		Count:   len(orders),
	}, nil
}

func (uc *DeliveryUseCase) GetOrderHistory(partnerID string, req *entities.GetOrderHistoryRequest) (*entities.GetOrderHistoryResponse, error) {
	deliveries, total, err := uc.deliveryRepo.GetOrderHistory(partnerID, req.Limit, req.Offset)
	if err != nil {
		return &entities.GetOrderHistoryResponse{
			Success: false,
		}, err
	}

	orders := make([]entities.DeliveryListItem, 0, len(deliveries))
	for _, d := range deliveries {
		orders = append(orders, entities.DeliveryListItem{
			ID:          d.DeliveryID,
			OrderID:     d.OrderID,
			Status:      d.Status,
			Address:     d.DeliveryAddress,
			Amount:      d.OrderAmount,
			DeliveryFee: d.DeliveryFee,
			ItemsCount:  d.ItemsCount,
			Distance:    d.Distance,
			CreatedAt:   d.CreatedAt,
		})
	}

	hasNext := (req.Offset + req.Limit) < total
	hasPrevious := req.Offset > 0

	return &entities.GetOrderHistoryResponse{
		Success:     true,
		Orders:      orders,
		Total:       total,
		Limit:       req.Limit,
		Offset:      req.Offset,
		HasNext:     hasNext,
		HasPrevious: hasPrevious,
	}, nil
}

func (uc *DeliveryUseCase) GetOrderDetails(deliveryID string) (*entities.GetOrderDetailsResponse, error) {
	delivery, err := uc.deliveryRepo.GetByID(deliveryID)
	if err != nil {
		return &entities.GetOrderDetailsResponse{
			Success: false,
		}, err
	}

	return &entities.GetOrderDetailsResponse{
		Success: true,
		Order:   delivery,
	}, nil
}

func (uc *DeliveryUseCase) AcceptOrder(deliveryID, partnerID string) (*entities.ResponseMessage, error) {
	delivery, err := uc.deliveryRepo.GetByID(deliveryID)
	if err != nil {
		return &entities.ResponseMessage{
			Success: false,
			Error:   "Order not found",
		}, err
	}

	if delivery.Status != "pending" {
		return &entities.ResponseMessage{
			Success: false,
			Message: "Order is not available",
		}, nil
	}

	if err := uc.deliveryRepo.AcceptOrder(deliveryID, partnerID); err != nil {
		return &entities.ResponseMessage{
			Success: false,
			Error:   "Failed to accept order",
		}, err
	}

	return &entities.ResponseMessage{
		Success: true,
		Message: "Order accepted successfully",
	}, nil
}

func (uc *DeliveryUseCase) UpdateOrderStatus(deliveryID, partnerID string, req *entities.UpdateOrderStatusRequest) (*entities.ResponseMessage, error) {
	delivery, err := uc.deliveryRepo.GetByID(deliveryID)
	if err != nil {
		return &entities.ResponseMessage{
			Success: false,
			Error:   "Order not found",
		}, err
	}

	if delivery.PartnerID != partnerID {
		return &entities.ResponseMessage{
			Success: false,
			Message: "Unauthorized",
		}, errors.New("unauthorized")
	}

	// Validate status transition
	validTransitions := map[string][]string{
		"pending":    {"picked_up"},
		"picked_up":  {"in_transit"},
		"in_transit": {"delivered"},
	}

	validNext, ok := validTransitions[delivery.Status]
	if !ok {
		return &entities.ResponseMessage{
			Success: false,
			Message: "Invalid current status",
		}, nil
	}

	isValid := false
	for _, status := range validNext {
		if status == req.Status {
			isValid = true
			break
		}
	}

	if !isValid {
		return &entities.ResponseMessage{
			Success: false,
			Message: "Invalid status transition",
		}, nil
	}

	if err := uc.deliveryRepo.UpdateStatus(deliveryID, req.Status); err != nil {
		return &entities.ResponseMessage{
			Success: false,
			Error:   "Failed to update status",
		}, err
	}

	// Update partner location
	if req.Latitude != 0 && req.Longitude != 0 {
		uc.partnerRepo.UpdateLocation(partnerID, req.Latitude, req.Longitude)
	}

	return &entities.ResponseMessage{
		Success: true,
		Message: "Status updated successfully",
	}, nil
}

func (uc *DeliveryUseCase) CompleteDelivery(deliveryID, partnerID string, req *entities.CompleteDeliveryRequest) (*entities.ResponseMessage, error) {
	delivery, err := uc.deliveryRepo.GetByID(deliveryID)
	if err != nil {
		return &entities.ResponseMessage{
			Success: false,
			Error:   "Order not found",
		}, err
	}

	if delivery.PartnerID != partnerID {
		return &entities.ResponseMessage{
			Success: false,
			Message: "Unauthorized",
		}, errors.New("unauthorized")
	}

	if delivery.Status != "in_transit" {
		return &entities.ResponseMessage{
			Success: false,
			Message: "Order must be in transit to complete",
		}, nil
	}

	if err := uc.deliveryRepo.CompleteDelivery(deliveryID, req.Notes); err != nil {
		return &entities.ResponseMessage{
			Success: false,
			Error:   "Failed to complete delivery",
		}, err
	}

	// Create earnings record
	bonus := 0
	if delivery.DeliveryFee > 100 {
		bonus = 10 // Bonus for high-value deliveries
	}

	earnings := &entities.Earnings{
		PartnerID:    partnerID,
		DeliveryID:   deliveryID,
		OrderID:      delivery.OrderID,
		Amount:       delivery.OrderAmount,
		DeliveryFee:  delivery.DeliveryFee,
		Bonus:        bonus,
		TotalEarning: delivery.DeliveryFee + bonus,
		EarnedAt:     time.Now(),
	}

	if err := uc.earningsRepo.Create(earnings); err != nil {
		// Log error but don't fail the completion
		// In production, use proper logging
	}

	// Update partner location
	if req.Latitude != 0 && req.Longitude != 0 {
		uc.partnerRepo.UpdateLocation(partnerID, req.Latitude, req.Longitude)
	}

	return &entities.ResponseMessage{
		Success: true,
		Message: "Delivery completed successfully",
	}, nil
}

