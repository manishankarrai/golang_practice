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
	flowerShopDB         = "golang_practice"
	flowerShopCollection = "flowershops"
)

// ErrFlowerShopNotFound is returned when a flower shop document does not exist.
var ErrFlowerShopNotFound = errors.New("flower shop not found")

// ErrFlowerAlreadyInShop is returned when a flower is already in a shop.
var ErrFlowerAlreadyInShop = errors.New("flower already exists in shop")

// ErrFlowerNotInShop is returned when a flower is not part of a shop.
var ErrFlowerNotInShop = errors.New("flower not found in shop")

func flowerShopColl() *mongo.Collection {
	return db.Database(flowerShopDB).Collection(flowerShopCollection)
}

// CreateFlowerShop inserts a new flower shop document with an empty flower list.
func CreateFlowerShop(ctx context.Context, shop models.FlowerShop) (models.FlowerShop, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	shop.ShopID = bson.NewObjectID()
	if shop.Flowers == nil {
		shop.Flowers = make([]models.Flower, 0)
	}
	if _, err := flowerShopColl().InsertOne(ctx, shop); err != nil {
		return models.FlowerShop{}, err
	}
	return shop, nil
}

// GetFlowerShops returns all flower shop documents.
func GetFlowerShops(ctx context.Context) ([]models.FlowerShop, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cursor, err := flowerShopColl().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	shops := make([]models.FlowerShop, 0)
	if err := cursor.All(ctx, &shops); err != nil {
		return nil, err
	}
	return shops, nil
}

// GetFlowerShopByID returns one flower shop by its shopid.
func GetFlowerShopByID(ctx context.Context, id string) (models.FlowerShop, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	shopID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.FlowerShop{}, ErrFlowerShopNotFound
	}

	var shop models.FlowerShop
	err = flowerShopColl().FindOne(ctx, bson.M{"shopid": shopID}).Decode(&shop)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.FlowerShop{}, ErrFlowerShopNotFound
		}
		return models.FlowerShop{}, err
	}
	if shop.Flowers == nil {
		shop.Flowers = make([]models.Flower, 0)
	}
	return shop, nil
}

// UpdateFlowerShop replaces the mutable shop fields without touching the flower list.
func UpdateFlowerShop(ctx context.Context, id string, shop models.FlowerShop) (models.FlowerShop, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	shopID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.FlowerShop{}, ErrFlowerShopNotFound
	}

	update := bson.M{"$set": bson.M{
		"name":     shop.Name,
		"owner":    shop.Owner,
		"location": shop.Location,
		"phone":    shop.Phone,
	}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.FlowerShop
	err = flowerShopColl().FindOneAndUpdate(ctx, bson.M{"shopid": shopID}, update, opts).Decode(&updated)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.FlowerShop{}, ErrFlowerShopNotFound
		}
		return models.FlowerShop{}, err
	}
	if updated.Flowers == nil {
		updated.Flowers = make([]models.Flower, 0)
	}
	return updated, nil
}

// DeleteFlowerShop removes a flower shop by its shopid.
func DeleteFlowerShop(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	shopID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return ErrFlowerShopNotFound
	}

	result, err := flowerShopColl().DeleteOne(ctx, bson.M{"shopid": shopID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrFlowerShopNotFound
	}
	return nil
}

// AddFlowerToShop pushes an existing flower into the shop's flowers array.
func AddFlowerToShop(ctx context.Context, shopID string, flowerID string) (models.FlowerShop, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	id, err := bson.ObjectIDFromHex(shopID)
	if err != nil {
		return models.FlowerShop{}, ErrFlowerShopNotFound
	}

	flower, err := GetFlowerByID(ctx, flowerID)
	if err != nil {
		return models.FlowerShop{}, err
	}

	filter := bson.M{
		"shopid":           id,
		"flowers.flowerid": bson.M{"$ne": flower.FlowerID},
	}
	update := bson.M{"$push": bson.M{"flowers": flower}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated models.FlowerShop
	err = flowerShopColl().FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			var existing models.FlowerShop
			if findErr := flowerShopColl().FindOne(ctx, bson.M{"shopid": id}).Decode(&existing); findErr != nil {
				return models.FlowerShop{}, ErrFlowerShopNotFound
			}
			return models.FlowerShop{}, ErrFlowerAlreadyInShop
		}
		return models.FlowerShop{}, err
	}
	return updated, nil
}

// RemoveFlowerFromShop pulls a flower out of the shop's flowers array.
func RemoveFlowerFromShop(ctx context.Context, shopID string, flowerID string) (models.FlowerShop, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	id, err := bson.ObjectIDFromHex(shopID)
	if err != nil {
		return models.FlowerShop{}, ErrFlowerShopNotFound
	}
	fid, err := bson.ObjectIDFromHex(flowerID)
	if err != nil {
		return models.FlowerShop{}, ErrFlowerNotInShop
	}

	filter := bson.M{"shopid": id, "flowers.flowerid": fid}
	update := bson.M{"$pull": bson.M{"flowers": bson.M{"flowerid": fid}}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated models.FlowerShop
	err = flowerShopColl().FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			var existing models.FlowerShop
			if findErr := flowerShopColl().FindOne(ctx, bson.M{"shopid": id}).Decode(&existing); findErr != nil {
				return models.FlowerShop{}, ErrFlowerShopNotFound
			}
			return models.FlowerShop{}, ErrFlowerNotInShop
		}
		return models.FlowerShop{}, err
	}
	return updated, nil
}
