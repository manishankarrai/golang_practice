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
	shopDB         = "golang_practice"
	shopCollection = "shops"
)

// ErrShopNotFound is returned when a shop document does not exist.
var ErrShopNotFound = errors.New("shop not found")

func shopColl() *mongo.Collection {
	return db.Database(shopDB).Collection(shopCollection)
}

// CreateShop inserts a new shop document into the MongoDB collection.
// Steps: set timeout, generate ID, insert document, return result or error.
func CreateShop(ctx context.Context, shop models.Shop) (models.Shop, error) {
	// Create a derived context with 10s timeout to prevent indefinite blocking on DB ops
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Generate a new unique MongoDB ObjectID for the shop (required for insert)
	shop.ID = bson.NewObjectID()

	// Perform the insert into the shops collection; ignore the inserted ID return value
	if _, err := shopColl().InsertOne(ctx, shop); err != nil {
		return models.Shop{}, err
	}

	// On success, return the shop (now containing the generated ID)
	return shop, nil
}

// GetShops returns all shop documents.
func GetShops(ctx context.Context) ([]models.Shop, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cursor, err := shopColl().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	shops := make([]models.Shop, 0)
	if err := cursor.All(ctx, &shops); err != nil {
		return nil, err
	}
	return shops, nil
}

// GetShopByID returns a single shop by its hex id.
func GetShopByID(ctx context.Context, id string) (models.Shop, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.Shop{}, ErrShopNotFound
	}

	var shop models.Shop
	err = shopColl().FindOne(ctx, bson.M{"_id": objID}).Decode(&shop)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Shop{}, ErrShopNotFound
		}
		return models.Shop{}, err
	}
	return shop, nil
}

// UpdateShop replaces the mutable fields of an existing shop.
func UpdateShop(ctx context.Context, id string, shop models.Shop) (models.Shop, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.Shop{}, ErrShopNotFound
	}

	update := bson.M{"$set": bson.M{
		"name":    shop.Name,
		"address": shop.Address,
		"city":    shop.City,
		"phone":   shop.Phone,
	}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Shop
	err = shopColl().FindOneAndUpdate(ctx, bson.M{"_id": objID}, update, opts).Decode(&updated)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Shop{}, ErrShopNotFound
		}
		return models.Shop{}, err
	}
	return updated, nil
}

// DeleteShop removes a shop by its hex id.
func DeleteShop(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return ErrShopNotFound
	}

	res, err := shopColl().DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrShopNotFound
	}
	return nil
}
