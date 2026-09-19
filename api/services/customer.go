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

var ErrCustomerNotFound = errors.New("customer not found")

func customerColl() *mongo.Collection {
	return db.Database("golang_practice").Collection("customers")
}

func customerObjectID(id string) (bson.ObjectID, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return bson.NilObjectID, ErrCustomerNotFound
	}
	return objectID, nil
}

func customerFields(customer models.Customer) bson.M {
	return bson.M{"name": customer.Name, "email": customer.Email, "phone": customer.Phone}
}

func CreateCustomer(ctx context.Context, customer models.Customer) (models.Customer, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	customer.ID = bson.NewObjectID()
	_, err := customerColl().InsertOne(ctx, customer)
	if err != nil {
		return models.Customer{}, err
	}
	return customer, nil
}

func GetCustomers(ctx context.Context) ([]models.Customer, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	cursor, err := customerColl().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	customers := make([]models.Customer, 0)
	if err := cursor.All(ctx, &customers); err != nil {
		return nil, err
	}
	return customers, nil
}

func GetCustomerByID(ctx context.Context, id string) (models.Customer, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	objectID, err := customerObjectID(id)
	if err != nil {
		return models.Customer{}, err
	}
	var customer models.Customer
	if err := customerColl().FindOne(ctx, bson.M{"_id": objectID}).Decode(&customer); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Customer{}, ErrCustomerNotFound
		}
		return models.Customer{}, err
	}
	return customer, nil
}

func UpdateCustomer(ctx context.Context, id string, customer models.Customer) (models.Customer, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	objectID, err := customerObjectID(id)
	if err != nil {
		return models.Customer{}, err
	}
	var updated models.Customer
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if err := customerColl().FindOneAndUpdate(ctx, bson.M{"_id": objectID}, bson.M{"$set": customerFields(customer)}, opts).Decode(&updated); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Customer{}, ErrCustomerNotFound
		}
		return models.Customer{}, err
	}
	return updated, nil
}

func UpdateCustomerLocation(ctx context.Context, id string, latitude, longitude float64) (models.Customer, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	objectID, err := customerObjectID(id)
	if err != nil {
		return models.Customer{}, err
	}
	var updated models.Customer
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	update := bson.M{"$set": bson.M{"latitude": latitude, "longitude": longitude, "updated_at": time.Now().UTC()}}
	if err := customerColl().FindOneAndUpdate(ctx, bson.M{"_id": objectID}, update, opts).Decode(&updated); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Customer{}, ErrCustomerNotFound
		}
		return models.Customer{}, err
	}
	return updated, nil
}

func DeleteCustomer(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	objectID, err := customerObjectID(id)
	if err != nil {
		return err
	}
	result, err := customerColl().DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrCustomerNotFound
	}
	return nil
}
