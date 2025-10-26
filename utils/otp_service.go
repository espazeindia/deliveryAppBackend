package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateOTP generates a 6-digit OTP
func GenerateOTP() int {
	rand.Seed(time.Now().UnixNano())
	return 100000 + rand.Intn(900000)
}

// SendOTP simulates sending OTP via SMS
// In production, integrate with SMS gateway like Twilio, AWS SNS, etc.
func SendOTP(phoneNumber string, otp int) error {
	// TODO: Integrate with actual SMS service
	fmt.Printf("📱 Sending OTP %d to %s\n", otp, phoneNumber)
	return nil
}

