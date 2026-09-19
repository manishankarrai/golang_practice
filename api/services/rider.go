package services

import (
	"context"
	"errors"
	"time"

	"practice/api/db"
	"practice/api/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const riderTimeout = 10 * time.Second

var ErrRiderNotFound = errors.New("rider not found")

func riderColl() *mongo.Collection {
	return db.Database("golang_practice").Collection("riders")
}

func riderObjectID(id string) (bson.ObjectID, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return bson.NilObjectID, ErrRiderNotFound
	}
	return objectID, nil
}

func riderFields(rider models.Rider) bson.M {
	return bson.M{"name": rider.Name, "email": rider.Email, "phone": rider.Phone}
}

func CreateRider(ctx context.Context, rider models.Rider) (models.Rider, error) {
	ctx, cancel := context.WithTimeout(ctx, riderTimeout)
	defer cancel()

	rider.ID = bson.NewObjectID()
	if _, err := riderColl().InsertOne(ctx, rider); err != nil {
		return models.Rider{}, err
	}
	return rider, nil
}

func GetRiders(ctx context.Context) ([]models.Rider, error) {
	ctx, cancel := context.WithTimeout(ctx, riderTimeout)
	defer cancel()

	cursor, err := riderColl().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	riders := make([]models.Rider, 0)
	if err := cursor.All(ctx, &riders); err != nil {
		return nil, err
	}
	return riders, nil
}

func GetRiderByID(ctx context.Context, id string) (models.Rider, error) {
	ctx, cancel := context.WithTimeout(ctx, riderTimeout)
	defer cancel()

	objectID, err := riderObjectID(id)
	if err != nil {
		return models.Rider{}, err
	}

	var rider models.Rider
	if err := riderColl().FindOne(ctx, bson.M{"_id": objectID}).Decode(&rider); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Rider{}, ErrRiderNotFound
		}
		return models.Rider{}, err
	}
	return rider, nil
}

func UpdateRider(ctx context.Context, id string, rider models.Rider) (models.Rider, error) {
	ctx, cancel := context.WithTimeout(ctx, riderTimeout)
	defer cancel()

	objectID, err := riderObjectID(id)
	if err != nil {
		return models.Rider{}, err
	}

	var updated models.Rider
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if err := riderColl().FindOneAndUpdate(ctx, bson.M{"_id": objectID}, bson.M{"$set": riderFields(rider)}, opts).Decode(&updated); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Rider{}, ErrRiderNotFound
		}
		return models.Rider{}, err
	}
	return updated, nil
}

func DeleteRider(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, riderTimeout)
	defer cancel()

	objectID, err := riderObjectID(id)
	if err != nil {
		return err
	}
	result, err := riderColl().DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrRiderNotFound
	}
	return nil
}
