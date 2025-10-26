package mongodb

import (
	"context"
	"deliveryAppBackend/config"
	"deliveryAppBackend/domain/entities"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EarningsMongoRepository struct {
	collection *mongo.Collection
}

func NewEarningsMongoRepository() *EarningsMongoRepository {
	return &EarningsMongoRepository{
		collection: config.GetCollection("earnings"),
	}
}

func (r *EarningsMongoRepository) Create(earnings *entities.Earnings) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	earnings.CreatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, earnings)
	if err != nil {
		return err
	}

	earnings.EarningsID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *EarningsMongoRepository) GetByPartnerID(partnerID string, period string) ([]entities.Earnings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
		"earnedAt": bson.M{
			"$gte": startDate,
		},
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "earnedAt", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var earnings []entities.Earnings
	if err = cursor.All(ctx, &earnings); err != nil {
		return nil, err
	}

	return earnings, nil
}

func (r *EarningsMongoRepository) GetHistory(partnerID string, limit, offset int) ([]entities.Earnings, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"partnerId": partnerID,
	}

	// Get total count
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	opts := options.Find().
		SetSort(bson.D{{Key: "earnedAt", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var earnings []entities.Earnings
	if err = cursor.All(ctx, &earnings); err != nil {
		return nil, 0, err
	}

	return earnings, int(total), nil
}

func (r *EarningsMongoRepository) GetTotalEarnings(partnerID string, period string) (int, error) {
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
		"earnedAt": bson.M{
			"$gte": startDate,
		},
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "total", Value: bson.D{{Key: "$sum", Value: "$totalEarning"}}},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []struct {
		Total int `bson:"total"`
	}
	if err = cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, nil
	}

	return result[0].Total, nil
}

func (r *EarningsMongoRepository) GetEarningsCount(partnerID string, period string) (int, error) {
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
		"earnedAt": bson.M{
			"$gte": startDate,
		},
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	return int(count), err
}

func (r *EarningsMongoRepository) GetAvgEarnings(partnerID string, period string) (int, error) {
	total, err := r.GetTotalEarnings(partnerID, period)
	if err != nil {
		return 0, err
	}

	count, err := r.GetEarningsCount(partnerID, period)
	if err != nil {
		return 0, err
	}

	if count == 0 {
		return 0, nil
	}

	return total / count, nil
}

