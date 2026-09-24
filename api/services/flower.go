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
	flowerDB         = "golang_practice"
	flowerCollection = "flowers"
)

// ErrFlowerNotFound is returned when a flower document does not exist.
var ErrFlowerNotFound = errors.New("flower not found")

func flowerColl() *mongo.Collection {
	return db.Database(flowerDB).Collection(flowerCollection)
}

// CreateFlower inserts a new flower document.
func CreateFlower(ctx context.Context, flower models.Flower) (models.Flower, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	flower.FlowerID = bson.NewObjectID()
	if _, err := flowerColl().InsertOne(ctx, flower); err != nil {
		return models.Flower{}, err
	}
	return flower, nil
}

// GetFlowers returns all flower documents.
func GetFlowers(ctx context.Context) ([]models.Flower, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cursor, err := flowerColl().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	flowers := make([]models.Flower, 0)
	if err := cursor.All(ctx, &flowers); err != nil {
		return nil, err
	}
	return flowers, nil
}

// GetFlowerByID returns one flower by its flowerid.
func GetFlowerByID(ctx context.Context, id string) (models.Flower, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	flowerID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.Flower{}, ErrFlowerNotFound
	}

	var flower models.Flower
	err = flowerColl().FindOne(ctx, bson.M{"flowerid": flowerID}).Decode(&flower)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Flower{}, ErrFlowerNotFound
		}
		return models.Flower{}, err
	}
	return flower, nil
}

// UpdateFlower replaces the mutable fields of an existing flower.
func UpdateFlower(ctx context.Context, id string, flower models.Flower) (models.Flower, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	flowerID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.Flower{}, ErrFlowerNotFound
	}

	update := bson.M{"$set": bson.M{
		"name":     flower.Name,
		"type":     flower.Type,
		"color":    flower.Color,
		"location": flower.Location,
		"pricing":  flower.Pricing,
		"details":  flower.Details,
	}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Flower
	err = flowerColl().FindOneAndUpdate(ctx, bson.M{"flowerid": flowerID}, update, opts).Decode(&updated)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Flower{}, ErrFlowerNotFound
		}
		return models.Flower{}, err
	}
	return updated, nil
}

// DeleteFlower removes a flower by its flowerid.
func DeleteFlower(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	flowerID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return ErrFlowerNotFound
	}

	result, err := flowerColl().DeleteOne(ctx, bson.M{"flowerid": flowerID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrFlowerNotFound
	}
	return nil
}
