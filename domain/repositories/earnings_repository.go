package repositories

import (
	"deliveryAppBackend/domain/entities"
)

type EarningsRepository interface {
	// Earnings Management
	Create(earnings *entities.Earnings) error
	GetByPartnerID(partnerID string, period string) ([]entities.Earnings, error)
	GetHistory(partnerID string, limit, offset int) ([]entities.Earnings, int, error)
	
	// Statistics
	GetTotalEarnings(partnerID string, period string) (int, error)
	GetEarningsCount(partnerID string, period string) (int, error)
	GetAvgEarnings(partnerID string, period string) (int, error)
}

