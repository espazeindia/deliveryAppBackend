package entities

import "time"

// Delivery represents a delivery assignment
type Delivery struct {
	DeliveryID       string    `json:"id" bson:"_id,omitempty"`
	OrderID          string    `json:"orderId" bson:"orderId"`
	PartnerID        string    `json:"partnerId" bson:"partnerId"`
	CustomerID       string    `json:"customerId" bson:"customerId"`
	CustomerName     string    `json:"customerName" bson:"customerName"`
	CustomerPhone    string    `json:"customerPhone" bson:"customerPhone"`
	WarehouseID      string    `json:"warehouseId" bson:"warehouseId"`
	Status           string    `json:"status" bson:"status"` // pending, picked_up, in_transit, delivered, cancelled
	PickupAddress    string    `json:"pickupAddress" bson:"pickupAddress"`
	DeliveryAddress  string    `json:"deliveryAddress" bson:"deliveryAddress"`
	PickupLatitude   float64   `json:"pickupLatitude" bson:"pickupLatitude"`
	PickupLongitude  float64   `json:"pickupLongitude" bson:"pickupLongitude"`
	DeliveryLatitude float64   `json:"deliveryLatitude" bson:"deliveryLatitude"`
	DeliveryLongitude float64  `json:"deliveryLongitude" bson:"deliveryLongitude"`
	Distance         float64   `json:"distance" bson:"distance"` // in km
	OrderAmount      int       `json:"orderAmount" bson:"orderAmount"`
	DeliveryFee      int       `json:"deliveryFee" bson:"deliveryFee"`
	ItemsCount       int       `json:"itemsCount" bson:"itemsCount"`
	Items            []OrderItem `json:"items" bson:"items"`
	AssignedAt       time.Time `json:"assignedAt" bson:"assignedAt"`
	PickedUpAt       *time.Time `json:"pickedUpAt,omitempty" bson:"pickedUpAt,omitempty"`
	InTransitAt      *time.Time `json:"inTransitAt,omitempty" bson:"inTransitAt,omitempty"`
	DeliveredAt      *time.Time `json:"deliveredAt,omitempty" bson:"deliveredAt,omitempty"`
	CreatedAt        time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt" bson:"updatedAt"`
	// Additional fields
	PaymentMethod    string    `json:"paymentMethod" bson:"paymentMethod"` // cod, online
	Notes            string    `json:"notes" bson:"notes"`
	CancellationReason string  `json:"cancellationReason,omitempty" bson:"cancellationReason,omitempty"`
}

type OrderItem struct {
	ProductID   string `json:"productId" bson:"productId"`
	Name        string `json:"name" bson:"name"`
	Quantity    int    `json:"quantity" bson:"quantity"`
	Price       int    `json:"price" bson:"price"`
	ImageURL    string `json:"imageUrl" bson:"imageUrl"`
}

// Requests and Responses

type GetActiveOrdersResponse struct {
	Success bool        `json:"success"`
	Orders  []DeliveryListItem `json:"orders"`
	Count   int         `json:"count"`
}

type DeliveryListItem struct {
	ID          string    `json:"id"`
	OrderID     string    `json:"orderId"`
	Status      string    `json:"status"`
	Address     string    `json:"address"`
	Amount      int       `json:"amount"`
	DeliveryFee int       `json:"deliveryFee"`
	ItemsCount  int       `json:"itemsCount"`
	Distance    float64   `json:"distance"`
	CreatedAt   time.Time `json:"createdAt"`
}

type GetOrderHistoryRequest struct {
	Limit  int `json:"limit" form:"limit" binding:"gte=1"`
	Offset int `json:"offset" form:"offset" binding:"gte=0"`
}

type GetOrderHistoryResponse struct {
	Success     bool               `json:"success"`
	Orders      []DeliveryListItem `json:"orders"`
	Total       int                `json:"total"`
	Limit       int                `json:"limit"`
	Offset      int                `json:"offset"`
	HasNext     bool               `json:"hasNext"`
	HasPrevious bool               `json:"hasPrevious"`
}

type GetOrderDetailsResponse struct {
	Success bool      `json:"success"`
	Order   *Delivery `json:"order"`
}

type UpdateOrderStatusRequest struct {
	Status    string  `json:"status" binding:"required"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type CompleteDeliveryRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Signature string  `json:"signature"`
	Notes     string  `json:"notes"`
}

