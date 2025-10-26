package entities

import "time"

// Earnings represents earnings for a delivery partner
type Earnings struct {
	EarningsID   string    `json:"id" bson:"_id,omitempty"`
	PartnerID    string    `json:"partnerId" bson:"partnerId"`
	DeliveryID   string    `json:"deliveryId" bson:"deliveryId"`
	OrderID      string    `json:"orderId" bson:"orderId"`
	Amount       int       `json:"amount" bson:"amount"`
	DeliveryFee  int       `json:"deliveryFee" bson:"deliveryFee"`
	Bonus        int       `json:"bonus" bson:"bonus"`
	TotalEarning int       `json:"totalEarning" bson:"totalEarning"`
	EarnedAt     time.Time `json:"earnedAt" bson:"earnedAt"`
	CreatedAt    time.Time `json:"createdAt" bson:"createdAt"`
}

// Requests and Responses

type GetEarningsRequest struct {
	Period string `json:"period" form:"period"` // today, week, month
}

type GetEarningsResponse struct {
	Success          bool    `json:"success"`
	TotalEarnings    int     `json:"totalEarnings"`
	DeliveriesCount  int     `json:"deliveriesCount"`
	AvgPerDelivery   int     `json:"avgPerDelivery"`
	BonusEarnings    int     `json:"bonusEarnings"`
	WeeklyCount      int     `json:"weeklyCount"`
	Period           string  `json:"period"`
}

type GetEarningsHistoryRequest struct {
	Limit  int `json:"limit" form:"limit" binding:"gte=1"`
	Offset int `json:"offset" form:"offset" binding:"gte=0"`
}

type EarningsHistoryItem struct {
	OrderID     string    `json:"orderId"`
	Amount      int       `json:"amount"`
	CompletedAt time.Time `json:"completedAt"`
}

type GetEarningsHistoryResponse struct {
	Success     bool                  `json:"success"`
	History     []EarningsHistoryItem `json:"history"`
	Total       int                   `json:"total"`
	Limit       int                   `json:"limit"`
	Offset      int                   `json:"offset"`
	HasNext     bool                  `json:"hasNext"`
	HasPrevious bool                  `json:"hasPrevious"`
}

