package entities

import "time"

// DeliveryPartner represents a delivery partner in the system
type DeliveryPartner struct {
	PartnerID          string    `json:"id" bson:"_id,omitempty"`
	Name               string    `json:"name" bson:"name"`
	PhoneNumber        string    `json:"phoneNumber" bson:"phoneNumber"`
	Email              string    `json:"email" bson:"email"`
	OTP                int       `json:"otp" bson:"otp"`
	NumberOfRetriesOTP int       `json:"numberOfRetriesOTP" bson:"numberOfRetriesOTP"`
	OTPGeneratedAt     time.Time `json:"otpGeneratedAt" bson:"otpGeneratedAt"`
	PIN                int       `json:"pin" bson:"pin"`
	NumberOfRetriesPIN int       `json:"numberOfRetriesPIN" bson:"numberOfRetriesPIN"`
	IsAvailable        bool      `json:"isAvailable" bson:"isAvailable"`
	IsVerified         bool      `json:"isVerified" bson:"isVerified"`
	Rating             float64   `json:"rating" bson:"rating"`
	TotalDeliveries    int       `json:"totalDeliveries" bson:"totalDeliveries"`
	LastLoginAt        time.Time `json:"lastLoginAt,omitempty" bson:"lastLoginAt,omitempty"`
	CreatedAt          time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt" bson:"updatedAt"`
	// Document verification
	AadharNumber      string `json:"aadharNumber" bson:"aadharNumber"`
	PanNumber         string `json:"panNumber" bson:"panNumber"`
	DrivingLicense    string `json:"drivingLicense" bson:"drivingLicense"`
	VehicleNumber     string `json:"vehicleNumber" bson:"vehicleNumber"`
	VehicleType       string `json:"vehicleType" bson:"vehicleType"` // bike, scooter, car
	BankAccountNumber string `json:"bankAccountNumber" bson:"bankAccountNumber"`
	IFSC              string `json:"ifsc" bson:"ifsc"`
	// Location
	CurrentLatitude  float64 `json:"currentLatitude" bson:"currentLatitude"`
	CurrentLongitude float64 `json:"currentLongitude" bson:"currentLongitude"`
	LastLocationAt   time.Time `json:"lastLocationAt" bson:"lastLocationAt"`
}

// Login Request/Response
type DeliveryPartnerLoginRequest struct {
	PhoneNumber string `json:"phoneNumber" binding:"required,min=10"`
	PIN         int    `json:"pin" binding:"required"`
}

type DeliveryPartnerLoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	User    *struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		PhoneNumber string `json:"phoneNumber"`
		IsAvailable bool   `json:"isAvailable"`
	} `json:"user,omitempty"`
	Error string `json:"error,omitempty"`
}

// OTP Request/Response
type RequestOTPRequest struct {
	PhoneNumber string `json:"phoneNumber" binding:"required,min=10"`
}

type VerifyOTPRequest struct {
	PhoneNumber string `json:"phoneNumber" binding:"required,min=10"`
	OTP         int    `json:"otp" binding:"required"`
}

type OTPResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	User    *struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		PhoneNumber string `json:"phoneNumber"`
		IsAvailable bool   `json:"isAvailable"`
	} `json:"user,omitempty"`
	Error string `json:"error,omitempty"`
}

// Profile Update
type UpdateProfileRequest struct {
	Name               string `json:"name"`
	Email              string `json:"email"`
	AadharNumber       string `json:"aadharNumber"`
	PanNumber          string `json:"panNumber"`
	DrivingLicense     string `json:"drivingLicense"`
	VehicleNumber      string `json:"vehicleNumber"`
	VehicleType        string `json:"vehicleType"`
	BankAccountNumber  string `json:"bankAccountNumber"`
	IFSC               string `json:"ifsc"`
}

// Location Update
type UpdateLocationRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

// Availability Toggle
type ToggleAvailabilityRequest struct {
	IsAvailable bool `json:"isAvailable"`
}

type ResponseMessage struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

