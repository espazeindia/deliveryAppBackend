package repositories

import (
	"deliveryAppBackend/domain/entities"
)

type DeliveryRepository interface {
	// Order Management
	GetActiveOrdersByPartner(partnerID string) ([]entities.Delivery, error)
	GetOrderHistory(partnerID string, limit, offset int) ([]entities.Delivery, int, error)
	GetByID(deliveryID string) (*entities.Delivery, error)
	GetByOrderID(orderID string) (*entities.Delivery, error)
	Create(delivery *entities.Delivery) error
	Update(delivery *entities.Delivery) error
	
	// Status Updates
	AcceptOrder(deliveryID, partnerID string) error
	UpdateStatus(deliveryID, status string) error
	CompleteDelivery(deliveryID string, notes string) error
	
	// Assignment
	AssignToPartner(orderID, partnerID string) error
	GetPendingOrders() ([]entities.Delivery, error)
	
	// Statistics
	GetDeliveriesCountByPartner(partnerID string, period string) (int, error)
}

