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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DeliveryMongoRepository struct {
	collection *mongo.Collection
}

func NewDeliveryMongoRepository() *DeliveryMongoRepository {
	return &DeliveryMongoRepository{
		collection: config.GetCollection("deliveries"),
	}
}

func (r *DeliveryMongoRepository) GetActiveOrdersByPartner(partnerID string) ([]entities.Delivery, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"partnerId": partnerID,
		"status":    bson.M{"$in": []string{"pending", "picked_up", "in_transit"}},
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var deliveries []entities.Delivery
	if err = cursor.All(ctx, &deliveries); err != nil {
		return nil, err
	}

	return deliveries, nil
}

func (r *DeliveryMongoRepository) GetOrderHistory(partnerID string, limit, offset int) ([]entities.Delivery, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"partnerId": partnerID,
		"status":    bson.M{"$in": []string{"delivered", "cancelled"}},
	}

	// Get total count
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	opts := options.Find().
		SetSort(bson.D{{Key: "deliveredAt", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var deliveries []entities.Delivery
	if err = cursor.All(ctx, &deliveries); err != nil {
		return nil, 0, err
	}

	return deliveries, int(total), nil
}

func (r *DeliveryMongoRepository) GetByID(deliveryID string) (*entities.Delivery, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(deliveryID)
	if err != nil {
		return nil, err
	}

	var delivery entities.Delivery
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&delivery)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("delivery not found")
		}
		return nil, err
	}

	return &delivery, nil
}

func (r *DeliveryMongoRepository) GetByOrderID(orderID string) (*entities.Delivery, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var delivery entities.Delivery
	err := r.collection.FindOne(ctx, bson.M{"orderId": orderID}).Decode(&delivery)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &delivery, nil
}

func (r *DeliveryMongoRepository) Create(delivery *entities.Delivery) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	delivery.CreatedAt = time.Now()
	delivery.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, delivery)
	if err != nil {
		return err
	}

	delivery.DeliveryID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *DeliveryMongoRepository) Update(delivery *entities.Delivery) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(delivery.DeliveryID)
	if err != nil {
		return err
	}

	delivery.UpdatedAt = time.Now()

	_, err = r.collection.ReplaceOne(ctx, bson.M{"_id": objectID}, delivery)
	return err
}

func (r *DeliveryMongoRepository) AcceptOrder(deliveryID, partnerID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(deliveryID)
	if err != nil {
		return err
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"partnerId":  partnerID,
			"status":     "picked_up",
			"pickedUpAt": now,
			"updatedAt":  now,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

func (r *DeliveryMongoRepository) UpdateStatus(deliveryID, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(deliveryID)
	if err != nil {
		return err
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":    status,
			"updatedAt": now,
		},
	}

	// Set appropriate timestamp based on status
	if status == "in_transit" {
		update["$set"].(bson.M)["inTransitAt"] = now
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

func (r *DeliveryMongoRepository) CompleteDelivery(deliveryID string, notes string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(deliveryID)
	if err != nil {
		return err
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":      "delivered",
			"deliveredAt": now,
			"notes":       notes,
			"updatedAt":   now,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

func (r *DeliveryMongoRepository) AssignToPartner(orderID, partnerID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"partnerId":  partnerID,
			"assignedAt": time.Now(),
			"updatedAt":  time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"orderId": orderID}, update)
	return err
}

func (r *DeliveryMongoRepository) GetPendingOrders() ([]entities.Delivery, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"status": "pending",
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "createdAt", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var deliveries []entities.Delivery
	if err = cursor.All(ctx, &deliveries); err != nil {
		return nil, err
	}

	return deliveries, nil
}

func (r *DeliveryMongoRepository) GetDeliveriesCountByPartner(partnerID string, period string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var startDate time.Time
	now := time.Now()

	switch period {
	case "today":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	case "week":
		startDate = now.AddDate(0, 0, -7)
	case "month":
		startDate = now.AddDate(0, -1, 0)
	default:
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	}

	filter := bson.M{
		"partnerId": partnerID,
		"status":    "delivered",
		"deliveredAt": bson.M{
			"$gte": startDate,
		},
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	return int(count), err
}

