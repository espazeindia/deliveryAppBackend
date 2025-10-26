package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	PartnerID   string `json:"partnerId"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	IsAvailable bool   `json:"isAvailable"`
	jwt.RegisteredClaims
}

func GenerateJWT(partnerID, name, phoneNumber string, isAvailable bool) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "delivery-secret-key-change-in-production"
	}

	claims := Claims{
		PartnerID:   partnerID,
		Name:        name,
		PhoneNumber: phoneNumber,
		IsAvailable: isAvailable,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Hour)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenString string) (*Claims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "delivery-secret-key-change-in-production"
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

