package mongodb

import (
	"context"
	"deliveryAppBackend/config"
	"deliveryAppBackend/domain/entities"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DeliveryPartnerMongoRepository struct {
	collection *mongo.Collection
}

func NewDeliveryPartnerMongoRepository() *DeliveryPartnerMongoRepository {
	return &DeliveryPartnerMongoRepository{
		collection: config.GetCollection("delivery_partners"),
	}
}

func (r *DeliveryPartnerMongoRepository) FindByPhoneNumber(phoneNumber string) (*entities.DeliveryPartner, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var partner entities.DeliveryPartner
	err := r.collection.FindOne(ctx, bson.M{"phoneNumber": phoneNumber}).Decode(&partner)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &partner, nil
}

func (r *DeliveryPartnerMongoRepository) FindByID(partnerID string) (*entities.DeliveryPartner, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(partnerID)
	if err != nil {
		return nil, err
	}

	var partner entities.DeliveryPartner
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&partner)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("partner not found")
		}
		return nil, err
	}

	return &partner, nil
}

func (r *DeliveryPartnerMongoRepository) Create(partner *entities.DeliveryPartner) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	partner.CreatedAt = time.Now()
	partner.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, partner)
	if err != nil {
		return err
	}

	partner.PartnerID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *DeliveryPartnerMongoRepository) Update(partner *entities.DeliveryPartner) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(partner.PartnerID)
	if err != nil {
		return err
	}

	partner.UpdatedAt = time.Now()

	_, err = r.collection.ReplaceOne(ctx, bson.M{"_id": objectID}, partner)
	return err
}

func (r *DeliveryPartnerMongoRepository) UpdateOTP(phoneNumber string, otp int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"otp":            otp,
			"otpGeneratedAt": time.Now(),
			"updatedAt":      time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"phoneNumber": phoneNumber}, update)
	return err
}

func (r *DeliveryPartnerMongoRepository) VerifyOTP(phoneNumber string, otp int) (*entities.DeliveryPartner, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var partner entities.DeliveryPartner
	err := r.collection.FindOne(ctx, bson.M{
		"phoneNumber": phoneNumber,
		"otp":         otp,
	}).Decode(&partner)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("invalid OTP")
		}
		return nil, err
	}

	// Check if OTP is expired (valid for 10 minutes)
	if time.Since(partner.OTPGeneratedAt) > 10*time.Minute {
		return nil, errors.New("OTP expired")
	}

	return &partner, nil
}

func (r *DeliveryPartnerMongoRepository) UpdateProfile(partnerID string, updates map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(partnerID)
	if err != nil {
		return err
	}

	updates["updatedAt"] = time.Now()

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updates})
	return err
}

func (r *DeliveryPartnerMongoRepository) UpdateLocation(partnerID string, latitude, longitude float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(partnerID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"currentLatitude":  latitude,
			"currentLongitude": longitude,
			"lastLocationAt":   time.Now(),
			"updatedAt":        time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

func (r *DeliveryPartnerMongoRepository) ToggleAvailability(partnerID string, isAvailable bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(partnerID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"isAvailable": isAvailable,
			"updatedAt":   time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

func (r *DeliveryPartnerMongoRepository) GetTotalDeliveries(partnerID string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(partnerID)
	if err != nil {
		return 0, err
	}

	var partner entities.DeliveryPartner
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&partner)
	if err != nil {
		return 0, err
	}

	return partner.TotalDeliveries, nil
}

func (r *DeliveryPartnerMongoRepository) UpdateRating(partnerID string, rating float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(partnerID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"rating":    rating,
			"updatedAt": time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

