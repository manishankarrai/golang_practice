package services

import (
	"context"
	"errors"
	"time"

	"practice/api/db"
	"practice/api/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
