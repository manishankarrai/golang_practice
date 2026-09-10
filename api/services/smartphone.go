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

const (
	smartphoneDB         = "golang_practice"
	smartphoneCollection = "smartphones"
)

// ErrSmartphoneNotFound is returned when a smartphone document does not exist.
var ErrSmartphoneNotFound = errors.New("smartphone not found")

func smartphoneColl() *mongo.Collection {
	return db.Database(smartphoneDB).Collection(smartphoneCollection)
}

// CreateSmartphone inserts a new smartphone document.
func CreateSmartphone(ctx context.Context, phone models.Smartphone) (models.Smartphone, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	phone.ID = bson.NewObjectID()
	if _, err := smartphoneColl().InsertOne(ctx, phone); err != nil {
		return models.Smartphone{}, err
	}
	return phone, nil
}

// GetSmartphones returns all smartphone documents, optionally filtered by category.
func GetSmartphones(ctx context.Context, category string) ([]models.Smartphone, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	filter := bson.M{}
	if category != "" {
		filter["category"] = category
	}

	cursor, err := smartphoneColl().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	phones := make([]models.Smartphone, 0)
	if err := cursor.All(ctx, &phones); err != nil {
		return nil, err
	}
	return phones, nil
}

// GetSmartphoneByID returns a single smartphone by its hex id.
func GetSmartphoneByID(ctx context.Context, id string) (models.Smartphone, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.Smartphone{}, ErrSmartphoneNotFound
	}

	var phone models.Smartphone
	err = smartphoneColl().FindOne(ctx, bson.M{"_id": objID}).Decode(&phone)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Smartphone{}, ErrSmartphoneNotFound
		}
		return models.Smartphone{}, err
	}
	return phone, nil
}

// UpdateSmartphone replaces the mutable fields of an existing smartphone.
func UpdateSmartphone(ctx context.Context, id string, phone models.Smartphone) (models.Smartphone, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.Smartphone{}, ErrSmartphoneNotFound
	}

	update := bson.M{"$set": bson.M{
		"name":     phone.Name,
		"brand":    phone.Brand,
		"category": phone.Category,
		"price":    phone.Price,
		"stock":    phone.Stock,
	}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Smartphone
	err = smartphoneColl().FindOneAndUpdate(ctx, bson.M{"_id": objID}, update, opts).Decode(&updated)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Smartphone{}, ErrSmartphoneNotFound
		}
		return models.Smartphone{}, err
	}
	return updated, nil
}

// DeleteSmartphone removes a smartphone by its hex id.
func DeleteSmartphone(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return ErrSmartphoneNotFound
	}

	res, err := smartphoneColl().DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrSmartphoneNotFound
	}
	return nil
}
