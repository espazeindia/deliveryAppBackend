package usecase

import (
	"deliveryAppBackend/domain/entities"
	"deliveryAppBackend/domain/repositories"
)

type ProfileUseCase struct {
	partnerRepo repositories.DeliveryPartnerRepository
}

func NewProfileUseCase(partnerRepo repositories.DeliveryPartnerRepository) *ProfileUseCase {
	return &ProfileUseCase{
		partnerRepo: partnerRepo,
	}
}

func (uc *ProfileUseCase) GetProfile(partnerID string) (*entities.DeliveryPartner, error) {
	return uc.partnerRepo.FindByID(partnerID)
}

func (uc *ProfileUseCase) UpdateProfile(partnerID string, req *entities.UpdateProfileRequest) (*entities.ResponseMessage, error) {
	updates := make(map[string]interface{})

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.AadharNumber != "" {
		updates["aadharNumber"] = req.AadharNumber
	}
	if req.PanNumber != "" {
		updates["panNumber"] = req.PanNumber
	}
	if req.DrivingLicense != "" {
		updates["drivingLicense"] = req.DrivingLicense
	}
	if req.VehicleNumber != "" {
		updates["vehicleNumber"] = req.VehicleNumber
	}
	if req.VehicleType != "" {
		updates["vehicleType"] = req.VehicleType
	}
	if req.BankAccountNumber != "" {
		updates["bankAccountNumber"] = req.BankAccountNumber
	}
	if req.IFSC != "" {
		updates["ifsc"] = req.IFSC
	}

	if err := uc.partnerRepo.UpdateProfile(partnerID, updates); err != nil {
		return &entities.ResponseMessage{
			Success: false,
			Error:   "Failed to update profile",
		}, err
	}

	return &entities.ResponseMessage{
		Success: true,
		Message: "Profile updated successfully",
	}, nil
}

func (uc *ProfileUseCase) UpdateLocation(partnerID string, req *entities.UpdateLocationRequest) (*entities.ResponseMessage, error) {
	if err := uc.partnerRepo.UpdateLocation(partnerID, req.Latitude, req.Longitude); err != nil {
		return &entities.ResponseMessage{
			Success: false,
			Error:   "Failed to update location",
		}, err
	}

	return &entities.ResponseMessage{
		Success: true,
		Message: "Location updated successfully",
	}, nil
}

func (uc *ProfileUseCase) ToggleAvailability(partnerID string, req *entities.ToggleAvailabilityRequest) (*entities.ResponseMessage, error) {
	if err := uc.partnerRepo.ToggleAvailability(partnerID, req.IsAvailable); err != nil {
		return &entities.ResponseMessage{
			Success: false,
			Error:   "Failed to update availability",
		}, err
	}

	status := "offline"
	if req.IsAvailable {
		status = "online"
	}

	return &entities.ResponseMessage{
		Success: true,
		Message: "You are now " + status,
	}, nil
}

