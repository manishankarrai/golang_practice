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

const registrationTimeout = 10 * time.Second

var ErrRegistrationNotFound = errors.New("registration not found")

func registrationColl() *mongo.Collection {
	return db.Database("golang_practice").Collection("registrations")
}

func registrationObjectID(id string) (bson.ObjectID, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return bson.NilObjectID, ErrRegistrationNotFound
	}
	return objID, nil
}

func registrationFields(registration models.Registration) bson.M {
	return bson.M{"vehicle_id": registration.VehicleID, "owner_name": registration.OwnerName, "owner_email": registration.OwnerEmail, "owner_phone": registration.OwnerPhone, "registration_number": registration.RegistrationNumber, "registration_date": registration.RegistrationDate, "expiry_date": registration.ExpiryDate, "owner_address": registration.OwnerAddress}
}

func CreateRegistration(ctx context.Context, registration models.Registration) (models.Registration, error) {
	ctx, cancel := context.WithTimeout(ctx, registrationTimeout)
	defer cancel()
	registration.ID = bson.NewObjectID()
	_, err := registrationColl().InsertOne(ctx, registration)
	if err != nil {
		return models.Registration{}, err
	}
	return registration, nil
}

func GetRegistrations(ctx context.Context) ([]models.Registration, error) {
	ctx, cancel := context.WithTimeout(ctx, registrationTimeout)
	defer cancel()
	cursor, err := registrationColl().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	registrations := make([]models.Registration, 0)
	if err := cursor.All(ctx, &registrations); err != nil {
		return nil, err
	}
	return registrations, nil
}

func GetRegistrationByID(ctx context.Context, id string) (models.Registration, error) {
	ctx, cancel := context.WithTimeout(ctx, registrationTimeout)
	defer cancel()
	objID, err := registrationObjectID(id)
	if err != nil {
		return models.Registration{}, err
	}
	var registration models.Registration
	if err := registrationColl().FindOne(ctx, bson.M{"_id": objID}).Decode(&registration); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Registration{}, ErrRegistrationNotFound
		}
		return models.Registration{}, err
	}
	return registration, nil
}

func UpdateRegistration(ctx context.Context, id string, registration models.Registration) (models.Registration, error) {
	ctx, cancel := context.WithTimeout(ctx, registrationTimeout)
	defer cancel()
	objID, err := registrationObjectID(id)
	if err != nil {
		return models.Registration{}, err
	}
	var updated models.Registration
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if err := registrationColl().FindOneAndUpdate(ctx, bson.M{"_id": objID}, bson.M{"$set": registrationFields(registration)}, opts).Decode(&updated); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Registration{}, ErrRegistrationNotFound
		}
		return models.Registration{}, err
	}
	return updated, nil
}

func DeleteRegistration(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, registrationTimeout)
	defer cancel()
	objID, err := registrationObjectID(id)
	if err != nil {
		return err
	}
	result, err := registrationColl().DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrRegistrationNotFound
	}
	return nil
}
