package repositories

import (
	"deliveryAppBackend/domain/entities"
)

type DeliveryPartnerRepository interface {
	// Authentication
	FindByPhoneNumber(phoneNumber string) (*entities.DeliveryPartner, error)
	FindByID(partnerID string) (*entities.DeliveryPartner, error)
	Create(partner *entities.DeliveryPartner) error
	Update(partner *entities.DeliveryPartner) error
	
	// OTP Management
	UpdateOTP(phoneNumber string, otp int) error
	VerifyOTP(phoneNumber string, otp int) (*entities.DeliveryPartner, error)
	
	// Profile Management
	UpdateProfile(partnerID string, updates map[string]interface{}) error
	UpdateLocation(partnerID string, latitude, longitude float64) error
	ToggleAvailability(partnerID string, isAvailable bool) error
	
	// Statistics
	GetTotalDeliveries(partnerID string) (int, error)
	UpdateRating(partnerID string, rating float64) error
}

