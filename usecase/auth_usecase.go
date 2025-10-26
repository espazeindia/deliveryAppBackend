package usecase

import (
	"deliveryAppBackend/domain/entities"
	"deliveryAppBackend/domain/repositories"
	"deliveryAppBackend/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	partnerRepo repositories.DeliveryPartnerRepository
}

func NewAuthUseCase(partnerRepo repositories.DeliveryPartnerRepository) *AuthUseCase {
	return &AuthUseCase{
		partnerRepo: partnerRepo,
	}
}

func (uc *AuthUseCase) Login(req *entities.DeliveryPartnerLoginRequest) (*entities.DeliveryPartnerLoginResponse, error) {
	partner, err := uc.partnerRepo.FindByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return &entities.DeliveryPartnerLoginResponse{
			Success: false,
			Error:   "Failed to find partner",
		}, err
	}

	if partner == nil {
		return &entities.DeliveryPartnerLoginResponse{
			Success: false,
			Message: "Partner not found. Please register first.",
		}, nil
	}

	// In a simple implementation, PIN is stored as integer
	// For production, you should hash the PIN
	if partner.PIN != req.PIN {
		return &entities.DeliveryPartnerLoginResponse{
			Success: false,
			Message: "Invalid PIN",
		}, nil
	}

	// Update last login using UpdateProfile to avoid _id issues
	updates := map[string]interface{}{
		"lastLoginAt": time.Now(),
	}
	if err := uc.partnerRepo.UpdateProfile(partner.PartnerID, updates); err != nil {
		// Log error but don't fail login
		// In production, use proper logging
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(partner.PartnerID, partner.Name, partner.PhoneNumber, partner.IsAvailable)
	if err != nil {
		return &entities.DeliveryPartnerLoginResponse{
			Success: false,
			Error:   "Failed to generate token",
		}, err
	}

	return &entities.DeliveryPartnerLoginResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
		User: &struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			PhoneNumber string `json:"phoneNumber"`
			IsAvailable bool   `json:"isAvailable"`
		}{
			ID:          partner.PartnerID,
			Name:        partner.Name,
			PhoneNumber: partner.PhoneNumber,
			IsAvailable: partner.IsAvailable,
		},
	}, nil
}

func (uc *AuthUseCase) RequestOTP(req *entities.RequestOTPRequest) (*entities.ResponseMessage, error) {
	partner, err := uc.partnerRepo.FindByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return &entities.ResponseMessage{
			Success: false,
			Error:   "Failed to process request",
		}, err
	}

	// Generate OTP
	otp := utils.GenerateOTP()

	if partner == nil {
		// Create new partner with OTP
		newPartner := &entities.DeliveryPartner{
			PhoneNumber:    req.PhoneNumber,
			OTP:            otp,
			OTPGeneratedAt: time.Now(),
			IsAvailable:    false,
			IsVerified:     false,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err := uc.partnerRepo.Create(newPartner); err != nil {
			return &entities.ResponseMessage{
				Success: false,
				Error:   "Failed to create partner",
			}, err
		}
	} else {
		// Update existing partner's OTP
		if err := uc.partnerRepo.UpdateOTP(req.PhoneNumber, otp); err != nil {
			return &entities.ResponseMessage{
				Success: false,
				Error:   "Failed to update OTP",
			}, err
		}
	}

	// Send OTP via SMS (in production, use actual SMS service)
	if err := utils.SendOTP(req.PhoneNumber, otp); err != nil {
		return &entities.ResponseMessage{
			Success: false,
			Error:   "Failed to send OTP",
		}, err
	}

	return &entities.ResponseMessage{
		Success: true,
		Message: "OTP sent successfully",
	}, nil
}

func (uc *AuthUseCase) VerifyOTP(req *entities.VerifyOTPRequest) (*entities.OTPResponse, error) {
	partner, err := uc.partnerRepo.VerifyOTP(req.PhoneNumber, req.OTP)
	if err != nil {
		return &entities.OTPResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Update last login using UpdateProfile to avoid _id issues
	updates := map[string]interface{}{
		"lastLoginAt": time.Now(),
	}
	if err := uc.partnerRepo.UpdateProfile(partner.PartnerID, updates); err != nil {
		// Log error but don't fail verification
		// In production, use proper logging
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(partner.PartnerID, partner.Name, partner.PhoneNumber, partner.IsAvailable)
	if err != nil {
		return &entities.OTPResponse{
			Success: false,
			Error:   "Failed to generate token",
		}, err
	}

	return &entities.OTPResponse{
		Success: true,
		Message: "OTP verified successfully",
		Token:   token,
		User: &struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			PhoneNumber string `json:"phoneNumber"`
			IsAvailable bool   `json:"isAvailable"`
		}{
			ID:          partner.PartnerID,
			Name:        partner.Name,
			PhoneNumber: partner.PhoneNumber,
			IsAvailable: partner.IsAvailable,
		},
	}, nil
}

func hashPIN(pin int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(string(rune(pin))), 14)
	return string(bytes), err
}

func checkPINHash(pin string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pin))
	return err == nil
}

