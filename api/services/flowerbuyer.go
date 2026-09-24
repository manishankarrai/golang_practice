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

var ErrFlowerBuyerNotFound = errors.New("flower buyer not found")

func flowerBuyerColl() *mongo.Collection {
	return db.Database("golang_practice").Collection("flowerbuyers")
}

func flowerBuyerID(id string) (bson.ObjectID, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return bson.NilObjectID, ErrFlowerBuyerNotFound
	}
	return objectID, nil
}

func CreateFlowerBuyer(ctx context.Context, buyer models.FlowerBuyer) (models.FlowerBuyer, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	buyer.BuyerID = bson.NewObjectID()
	if _, err := flowerBuyerColl().InsertOne(ctx, buyer); err != nil {
		return models.FlowerBuyer{}, err
	}
	return buyer, nil
}

func GetFlowerBuyers(ctx context.Context) ([]models.FlowerBuyer, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	cursor, err := flowerBuyerColl().Find(ctx, bson.M{})
	if err != nil { return nil, err }
	defer cursor.Close(ctx)
	buyers := make([]models.FlowerBuyer, 0)
	if err := cursor.All(ctx, &buyers); err != nil { return nil, err }
	return buyers, nil
}

func GetFlowerBuyerByID(ctx context.Context, id string) (models.FlowerBuyer, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	buyerID, err := flowerBuyerID(id)
	if err != nil { return models.FlowerBuyer{}, err }
	var buyer models.FlowerBuyer
	if err := flowerBuyerColl().FindOne(ctx, bson.M{"buyerid": buyerID}).Decode(&buyer); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) { return models.FlowerBuyer{}, ErrFlowerBuyerNotFound }
		return models.FlowerBuyer{}, err
	}
	return buyer, nil
}

func UpdateFlowerBuyer(ctx context.Context, id string, buyer models.FlowerBuyer) (models.FlowerBuyer, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	buyerID, err := flowerBuyerID(id)
	if err != nil { return models.FlowerBuyer{}, err }
	update := bson.M{"$set": bson.M{"name": buyer.Name, "email": buyer.Email, "phone": buyer.Phone, "address": buyer.Address}}
	var updated models.FlowerBuyer
	err = flowerBuyerColl().FindOneAndUpdate(ctx, bson.M{"buyerid": buyerID}, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updated)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) { return models.FlowerBuyer{}, ErrFlowerBuyerNotFound }
		return models.FlowerBuyer{}, err
	}
	return updated, nil
}

func DeleteFlowerBuyer(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	buyerID, err := flowerBuyerID(id)
	if err != nil { return err }
	result, err := flowerBuyerColl().DeleteOne(ctx, bson.M{"buyerid": buyerID})
	if err != nil { return err }
	if result.DeletedCount == 0 { return ErrFlowerBuyerNotFound }
	return nil
}
