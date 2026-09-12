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
	countryDB         = "golang_practice"
	countryCollection = "countries"
)

// ErrCountryNotFound is returned when a country document does not exist.
var ErrCountryNotFound = errors.New("country not found")

func countryColl() *mongo.Collection {
	return db.Database(countryDB).Collection(countryCollection)
}

// CreateCountry inserts a new country document.
func CreateCountry(ctx context.Context, country models.Country) (models.Country, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	country.ID = bson.NewObjectID()
	if _, err := countryColl().InsertOne(ctx, country); err != nil {
		return models.Country{}, err
	}
	return country, nil
}

// GetCountries returns all country documents.
func GetCountries(ctx context.Context) ([]models.Country, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cursor, err := countryColl().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	countries := make([]models.Country, 0)
	if err := cursor.All(ctx, &countries); err != nil {
		return nil, err
	}
	return countries, nil
}

// GetCountryByID returns a single country by its hex id.
func GetCountryByID(ctx context.Context, id string) (models.Country, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.Country{}, ErrCountryNotFound
	}

	var country models.Country
	err = countryColl().FindOne(ctx, bson.M{"_id": objID}).Decode(&country)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Country{}, ErrCountryNotFound
		}
		return models.Country{}, err
	}
	return country, nil
}

// UpdateCountry replaces the mutable fields of an existing country.
func UpdateCountry(ctx context.Context, id string, country models.Country) (models.Country, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.Country{}, ErrCountryNotFound
	}

	update := bson.M{"$set": bson.M{
		"name": country.Name,
		"code": country.Code,
	}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Country
	err = countryColl().FindOneAndUpdate(ctx, bson.M{"_id": objID}, update, opts).Decode(&updated)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Country{}, ErrCountryNotFound
		}
		return models.Country{}, err
	}
	return updated, nil
}

// DeleteCountry removes a country by its hex id.
func DeleteCountry(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return ErrCountryNotFound
	}

	res, err := countryColl().DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrCountryNotFound
	}
	return nil
}
