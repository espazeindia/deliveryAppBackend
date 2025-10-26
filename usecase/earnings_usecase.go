package usecase

import (
	"deliveryAppBackend/domain/entities"
	"deliveryAppBackend/domain/repositories"
)

type EarningsUseCase struct {
	earningsRepo repositories.EarningsRepository
	deliveryRepo repositories.DeliveryRepository
}

func NewEarningsUseCase(
	earningsRepo repositories.EarningsRepository,
	deliveryRepo repositories.DeliveryRepository,
) *EarningsUseCase {
	return &EarningsUseCase{
		earningsRepo: earningsRepo,
		deliveryRepo: deliveryRepo,
	}
}

func (uc *EarningsUseCase) GetEarnings(partnerID string, period string) (*entities.GetEarningsResponse, error) {
	if period == "" {
		period = "week"
	}

	totalEarnings, err := uc.earningsRepo.GetTotalEarnings(partnerID, period)
	if err != nil {
		return &entities.GetEarningsResponse{
			Success: false,
		}, err
	}

	deliveriesCount, err := uc.earningsRepo.GetEarningsCount(partnerID, period)
	if err != nil {
		return &entities.GetEarningsResponse{
			Success: false,
		}, err
	}

	avgPerDelivery, err := uc.earningsRepo.GetAvgEarnings(partnerID, period)
	if err != nil {
		return &entities.GetEarningsResponse{
			Success: false,
		}, err
	}

	// Get weekly count for additional stats
	weeklyCount, _ := uc.deliveryRepo.GetDeliveriesCountByPartner(partnerID, "week")

	// Calculate bonus earnings (15% of total)
	bonusEarnings := int(float64(totalEarnings) * 0.15)

	return &entities.GetEarningsResponse{
		Success:         true,
		TotalEarnings:   totalEarnings,
		DeliveriesCount: deliveriesCount,
		AvgPerDelivery:  avgPerDelivery,
		BonusEarnings:   bonusEarnings,
		WeeklyCount:     weeklyCount,
		Period:          period,
	}, nil
}

func (uc *EarningsUseCase) GetEarningsHistory(partnerID string, req *entities.GetEarningsHistoryRequest) (*entities.GetEarningsHistoryResponse, error) {
	earnings, total, err := uc.earningsRepo.GetHistory(partnerID, req.Limit, req.Offset)
	if err != nil {
		return &entities.GetEarningsHistoryResponse{
			Success: false,
		}, err
	}

	history := make([]entities.EarningsHistoryItem, 0, len(earnings))
	for _, e := range earnings {
		history = append(history, entities.EarningsHistoryItem{
			OrderID:     e.OrderID,
			Amount:      e.TotalEarning,
			CompletedAt: e.EarnedAt,
		})
	}

	hasNext := (req.Offset + req.Limit) < total
	hasPrevious := req.Offset > 0

	return &entities.GetEarningsHistoryResponse{
		Success:     true,
		History:     history,
		Total:       total,
		Limit:       req.Limit,
		Offset:      req.Offset,
		HasNext:     hasNext,
		HasPrevious: hasPrevious,
	}, nil
}

