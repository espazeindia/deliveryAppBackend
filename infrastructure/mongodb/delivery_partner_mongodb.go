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

	var result bson.M
	err := r.collection.FindOne(ctx, bson.M{"phoneNumber": phoneNumber}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	// Extract and remove _id before conversion
	var partnerID string
	if id, ok := result["_id"].(primitive.ObjectID); ok {
		partnerID = id.Hex()
		delete(result, "_id") // Remove _id from the map
	}

	// Convert to entity
	var partner entities.DeliveryPartner
	bsonBytes, _ := bson.Marshal(result)
	bson.Unmarshal(bsonBytes, &partner)
	
	// Set PartnerID
	partner.PartnerID = partnerID

	return &partner, nil
}

func (r *DeliveryPartnerMongoRepository) FindByID(partnerID string) (*entities.DeliveryPartner, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(partnerID)
	if err != nil {
		return nil, err
	}

	var result bson.M
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("partner not found")
		}
		return nil, err
	}

	// Remove _id from the map before conversion
	delete(result, "_id")

	// Convert to entity
	var partner entities.DeliveryPartner
	bsonBytes, _ := bson.Marshal(result)
	bson.Unmarshal(bsonBytes, &partner)
	partner.PartnerID = objectID.Hex()

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

	// Use UpdateOne with $set to only update specific fields
	update := bson.M{
		"$set": bson.M{
			"lastLoginAt": time.Now(),
			"updatedAt":   time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
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

	var result bson.M
	err := r.collection.FindOne(ctx, bson.M{
		"phoneNumber": phoneNumber,
		"otp":         otp,
	}).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("invalid OTP")
		}
		return nil, err
	}

	// Extract and remove _id before conversion
	var partnerID string
	if id, ok := result["_id"].(primitive.ObjectID); ok {
		partnerID = id.Hex()
		delete(result, "_id") // Remove _id from the map
	}

	// Convert to entity
	var partner entities.DeliveryPartner
	bsonBytes, _ := bson.Marshal(result)
	bson.Unmarshal(bsonBytes, &partner)
	
	// Set PartnerID
	partner.PartnerID = partnerID

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

	// Make sure _id is not in updates
	delete(updates, "_id")
	delete(updates, "id")
	
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

