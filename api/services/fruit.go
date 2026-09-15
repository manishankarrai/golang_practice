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
	fruitDB         = "golang_practice"
	fruitCollection = "fruits"
)

// ErrFruitNotFound is returned when a fruit document does not exist.
var ErrFruitNotFound = errors.New("fruit not found")

func fruitColl() *mongo.Collection {
	return db.Database(fruitDB).Collection(fruitCollection)
}

// CreateFruit inserts a new fruit document.
func CreateFruit(ctx context.Context, fruit models.Fruit) (models.Fruit, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	fruit.FruitID = bson.NewObjectID()
	if _, err := fruitColl().InsertOne(ctx, fruit); err != nil {
		return models.Fruit{}, err
	}
	return fruit, nil
}

// GetFruits returns all fruit documents.
func GetFruits(ctx context.Context) ([]models.Fruit, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cursor, err := fruitColl().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	fruits := make([]models.Fruit, 0)
	if err := cursor.All(ctx, &fruits); err != nil {
		return nil, err
	}
	return fruits, nil
}

// GetFruitByID returns one fruit by its fruitid.
func GetFruitByID(ctx context.Context, id string) (models.Fruit, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	fruitID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.Fruit{}, ErrFruitNotFound
	}

	var fruit models.Fruit
	err = fruitColl().FindOne(ctx, bson.M{"fruitid": fruitID}).Decode(&fruit)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Fruit{}, ErrFruitNotFound
		}
		return models.Fruit{}, err
	}
	return fruit, nil
}

// UpdateFruit replaces the mutable fields of an existing fruit.
func UpdateFruit(ctx context.Context, id string, fruit models.Fruit) (models.Fruit, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	fruitID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.Fruit{}, ErrFruitNotFound
	}

	update := bson.M{"$set": bson.M{
		"name":     fruit.Name,
		"type":     fruit.Type,
		"location": fruit.Location,
		"pricing":  fruit.Pricing,
		"details":  fruit.Details,
	}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Fruit
	err = fruitColl().FindOneAndUpdate(ctx, bson.M{"fruitid": fruitID}, update, opts).Decode(&updated)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Fruit{}, ErrFruitNotFound
		}
		return models.Fruit{}, err
	}
	return updated, nil
}

// DeleteFruit removes a fruit by its fruitid.
func DeleteFruit(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	fruitID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return ErrFruitNotFound
	}

	result, err := fruitColl().DeleteOne(ctx, bson.M{"fruitid": fruitID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrFruitNotFound
	}
	return nil
}
