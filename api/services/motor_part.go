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

const motorPartTimeout = 10 * time.Second

var ErrMotorPartNotFound = errors.New("motor part not found")

func motorPartColl() *mongo.Collection {
	return db.Database("golang_practice").Collection("motor_parts")
}

func motorPartObjectID(id string) (bson.ObjectID, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return bson.NilObjectID, ErrMotorPartNotFound
	}
	return objID, nil
}

func motorPartFields(part models.MotorPart) bson.M {
	return bson.M{
		"name":        part.Name,
		"description": part.Description,
		"price":       part.Price,
	}
}

func CreateMotorPart(ctx context.Context, part models.MotorPart) (models.MotorPart, error) {
	ctx, cancel := context.WithTimeout(ctx, motorPartTimeout)
	defer cancel()
	part.ID = bson.NewObjectID()
	_, err := motorPartColl().InsertOne(ctx, part)
	if err != nil {
		return models.MotorPart{}, err
	}
	return part, nil
}

func GetMotorParts(ctx context.Context) ([]models.MotorPart, error) {
	ctx, cancel := context.WithTimeout(ctx, motorPartTimeout)
	defer cancel()
	cursor, err := motorPartColl().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	parts := make([]models.MotorPart, 0)
	if err := cursor.All(ctx, &parts); err != nil {
		return nil, err
	}
	return parts, nil
}

func GetMotorPartByID(ctx context.Context, id string) (models.MotorPart, error) {
	ctx, cancel := context.WithTimeout(ctx, motorPartTimeout)
	defer cancel()
	objID, err := motorPartObjectID(id)
	if err != nil {
		return models.MotorPart{}, err
	}
	var part models.MotorPart
	if err := motorPartColl().FindOne(ctx, bson.M{"_id": objID}).Decode(&part); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.MotorPart{}, ErrMotorPartNotFound
		}
		return models.MotorPart{}, err
	}
	return part, nil
}

func UpdateMotorPart(ctx context.Context, id string, part models.MotorPart) (models.MotorPart, error) {
	ctx, cancel := context.WithTimeout(ctx, motorPartTimeout)
	defer cancel()
	objID, err := motorPartObjectID(id)
	if err != nil {
		return models.MotorPart{}, err
	}
	var updated models.MotorPart
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if err := motorPartColl().FindOneAndUpdate(ctx, bson.M{"_id": objID}, bson.M{"$set": motorPartFields(part)}, opts).Decode(&updated); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.MotorPart{}, ErrMotorPartNotFound
		}
		return models.MotorPart{}, err
	}
	return updated, nil
}

func DeleteMotorPart(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, motorPartTimeout)
	defer cancel()
	objID, err := motorPartObjectID(id)
	if err != nil {
		return err
	}
	result, err := motorPartColl().DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrMotorPartNotFound
	}
	return nil
}
